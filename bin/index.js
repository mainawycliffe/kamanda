#!/usr/bin/env node


"use strict";

var _typeof = typeof Symbol === "function" && typeof Symbol.iterator === "symbol" ? function (obj) { return typeof obj; } : function (obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj; };

var request = require("request"),
    path = require("path"),
    unzipper = require("unzipper"),
    mkdirp = require("mkdirp"),
    fs = require("fs"),
    exec = require("child_process").exec;

// Mapping from Node's `process.arch` to Golang's `$GOARCH`
var ARCH_MAPPING = {
  ia32: "386",
  x64: "amd64",
  arm: "arm"
};

// Mapping between Node's `process.platform` to Golang's
var PLATFORM_MAPPING = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
  freebsd: "freebsd"
};

function getInstallationPath(callback) {
  // `$npm_execpath bin` will output the path where binary files should be installed
  // using whichever package manager is current
  var execPath = process.env.npm_execpath;
  var packageManager = execPath && execPath.includes("yarn") ? "yarn global" : "npm";
  exec(packageManager + " bin", function (err, stdout, stderr) {
    var dir = null;
    if (err || stderr && !stderr.includes("No license field") || !stdout || stdout.length === 0) {
      // We couldn't infer path from `npm bin`. Let's try to get it from
      // Environment variables set by NPM when it runs.
      // npm_config_prefix points to NPM's installation directory where `bin` folder is available
      // Ex: /Users/foo/.nvm/versions/node/v4.3.0
      var env = process.env;
      if (env && env.npm_config_prefix) {
        dir = path.join(env.npm_config_prefix, "bin");
      }
    } else {
      dir = stdout.trim();
    }

    mkdirp.sync(dir);

    callback(null, dir);
  });
}

function verifyAndPlaceBinary(binName, binPath, callback) {
  if (!fs.existsSync(path.join(binPath, binName))) return callback("Downloaded binary does not contain the binary specified in configuration - " + binName);

  getInstallationPath(function (err, installationPath) {
    if (err) return callback("Error getting binary installation path from `npm bin`");

    // Move the binary file
    fs.renameSync(path.join(binPath, binName), path.join(installationPath, binName));

    callback(null);
  });
}

function validateConfiguration(packageJson) {
  if (!packageJson.version) {
    return "'version' property must be specified";
  }

  if (!packageJson.goBinary || _typeof(packageJson.goBinary) !== "object") {
    return "'goBinary' property must be defined and be an object";
  }

  if (!packageJson.goBinary.name) {
    return "'name' property is necessary";
  }

  if (!packageJson.goBinary.path) {
    return "'path' property is necessary";
  }

  if (!packageJson.goBinary.url) {
    return "'url' property is required";
  }

  // if (!packageJson.bin || typeof(packageJson.bin) !== "object") {
  //     return "'bin' property of package.json must be defined and be an object";
  // }
}

function parsePackageJson() {
  if (!(process.arch in ARCH_MAPPING)) {
    console.error("Installation is not supported for this architecture: " + process.arch);
    return;
  }

  if (!(process.platform in PLATFORM_MAPPING)) {
    console.error("Installation is not supported for this platform: " + process.platform);
    return;
  }

  var packageJsonPath = path.join(".", "package.json");
  if (!fs.existsSync(packageJsonPath)) {
    console.error("Unable to find package.json. " + "Please run this script at root of the package you want to be installed");
    return;
  }

  var packageJson = JSON.parse(fs.readFileSync(packageJsonPath));
  var error = validateConfiguration(packageJson);
  if (error && error.length > 0) {
    console.error("Invalid package.json: " + error);
    return;
  }

  // We have validated the config. It exists in all its glory
  var binName = packageJson.goBinary.name;
  var binPath = packageJson.goBinary.path;
  var url = packageJson.goBinary.url;
  var version = packageJson.version;
  if (version[0] === "v") version = version.substr(1); // strip the 'v' if necessary v0.0.1 => 0.0.1

  // Binary name on Windows has .exe suffix
  if (process.platform === "win32") {
    binName += ".exe";
  }

  // Interpolate variables in URL, if necessary
  url = url.replace(/{{arch}}/g, ARCH_MAPPING[process.arch]);
  url = url.replace(/{{platform}}/g, PLATFORM_MAPPING[process.platform]);
  url = url.replace(/{{version}}/g, version);
  url = url.replace(/{{bin_name}}/g, binName);

  return {
    binName: binName,
    binPath: binPath,
    url: url,
    version: version
  };
}

/**
 * Reads the configuration from application's package.json,
 * validates properties, downloads the binary, untars, and stores at
 * ./bin in the package's root. NPM already has support to install binary files
 * specific locations when invoked with "npm install -g"
 *
 *  See: https://docs.npmjs.com/files/package.json#bin
 */
function install(callback) {
  var opts = parsePackageJson();
  if (!opts) {
    return callback("Invalid inputs");
  }
  mkdirp.sync(opts.binPath);
  console.log("Downloading from URL: " + opts.url);
  var req = request({ uri: opts.url });
  req.on("error", callback.bind(null, "Error downloading from URL: " + opts.url));
  req.on("response", function (res) {
    if (res.statusCode !== 200) {
      return callback("Error downloading binary. HTTP Status Code: " + res.statusCode);
    }
    req.pipe(unzipper.Extract({ path: opts.binPath })).on("error", callback)
    // First we will Un-GZip, then we will untar. So once untar is completed,
    // binary is downloaded into `binPath`. Verify the binary and call it good
    .on("close", verifyAndPlaceBinary.bind(null, opts.binName, opts.binPath, callback));
  });
}

function uninstall(callback) {
  var opts = parsePackageJson();
  getInstallationPath(function (err, installationPath) {
    if (err) callback("Error finding Kamanda installation directory");

    try {
      fs.unlinkSync(path.join(installationPath, opts.binName));
    } catch (ex) {
      // Ignore errors when deleting the file.
    }

    return callback(null);
  });
}

// Parse command line arguments and call the right method
var actions = {
  install: install,
  uninstall: uninstall
};

var argv = process.argv;
if (argv && argv.length > 2) {
  var cmd = process.argv[2];
  if (!actions[cmd]) {
    console.log("Invalid command to go-npm. `install` and `uninstall` are the only supported commands");
    process.exit(1);
  }

  actions[cmd](function (err) {
    if (err) {
      console.error(err);
      process.exit(1);
    } else {
      process.exit(0);
    }
  });
}

// this code is thanks to:
// https://github.com/jumoel/go-npm/blob/feat/yarn-support/src/index.js and
// https://github.com/sanathkr/go-npm