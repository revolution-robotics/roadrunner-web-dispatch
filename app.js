#!/usr/bin/env node
import Koa from 'koa';
import mount from 'koa-mount';
import serve from 'koa-static'
import { exec, execFile } from 'child_process';
import { promisify } from 'util';
import parseArgs from 'minimist';
import ini from 'ini';
import fs from 'fs';

const app = new Koa();
const execPromise = promisify(exec);
const execFilePromise = promisify(execFile);

const pgm = process.argv[1].replace(/^.*\//, '');
const argv = parseArgs(process.argv.slice(2));

// runCmd: Returns the output of a shell command, cmd.
//
// "Never pass unsanitized user input to this function. Any input
//  containing shell metacharacters may be used to trigger
//  arbitrary command execution."
//
async function runCmd(cmd) {
    const { stdout, stderr } = await execPromise(cmd, { shell: '/bin/bash' });
    return stdout ? stdout : stderr;
}

// runFile: Returns the output of filePath invoked  with arguments, argList.
//
// If shell globbing, variable expansion, and I/O redirection are not
// needed, runFile should be prefered over runCmd.
//
async function runFile(filePath, argList = []) {
    const { stdout, stderr } = await execFilePromise(filePath, argList);
    return stdout ? stdout : stderr;
}

// As promise-based asynchronous functions, `runCmd' and `RunFile' can
// be invoked as either then-catch blocks or try-await-catch blocks.
// `await' can only be called outside an asynchronous function, as in
// the second example below, with JavaScript modules. To enable
// Javascript modules, include in `package.json' the declaration
// "type": "module" and use the `import' keyword instead of `require'.
//
// runFile('ls', ['-l', '/usr'])
//     .then(console.log)
//     .catch(error => console.log(`error: ${error}`));

// try {
//     console.log(await runFile('nmcli', ['--help']));
// } catch(err) {
//     console.error(`error: ${err}`);
// }


if (argv.help || argv.h) {
    console.log(`Usage: ${pgm} OPTIONS`);
    console.log(`OPTIONS:
    --config=PATH    Config file
    --cgi=PATH       CGI path (default: '/usr/bin/status.py')
    --port=N         Port to listen on (default: 80)
    --uri=TRIGGER    URI of CGI trigger (default: '/status.json')
    --www=PATH       HTML directory (default: '/var/www/html')`);
    process.exit(0);
}

if (argv.config) {
    if (! fs.existsSync(argv.config)) {
        console.error(`${pgm}: ${argv.config}: No such file or directory`);
        process.exit(1);
    }

    const config = ini.parse(fs.readFileSync(argv.config, 'utf-8'));

    // Command-line arguments override settings in config file, so
    // assign values to argv only if they're missing.
    for (const [key, value] of Object.entries(config)) {
        if (! argv[key]) {
            argv[key] = value;
        }
    }
}

async function runCGI(ctx, next) {
    await next();
    ctx.body = await runFile(argv.cgi);
}

async function redirectCockpit(ctx, next) {

    // Remove any port suffix from origin (e.g., :8080)
    const origin = ctx.request.origin.replace(/:[0-9]{1,5}$/, '');
    await next();
    ctx.redirect(`${origin}:9090`);
}

console.log(`Listening on port: ${argv.port}`);

app.use(mount(argv.uri, runCGI));
app.use(mount('/', serve(argv.www)));
app.use(mount('/cockpit', redirectCockpit));
app.listen(arg.port);
