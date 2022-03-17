/*
  @license
	Rollup.js v2.67.0
	Wed, 02 Feb 2022 06:19:08 GMT - commit 61d772ba59be605e53e72613ae51cee4b72aee64


	https://github.com/rollup/rollup

	Released under the MIT License.
*/
'use strict';

require('fs');
require('path');
require('url');
const loadConfigFile_js = require('./shared/loadConfigFile.js');
require('./shared/rollup.js');
require('./shared/mergeOptions.js');
require('process');
require('tty');
require('perf_hooks');
require('crypto');
require('events');



module.exports = loadConfigFile_js.loadAndParseConfigFile;
//# sourceMappingURL=loadConfigFile.js.map
