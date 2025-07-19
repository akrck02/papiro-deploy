import { spawn } from "node:child_process";

function chooseBinary() {
	return `main-linux-amd64`;
}

const binary = chooseBinary();
const mainScript = `./${binary}`;
const action = spawn(mainScript);
action.stdout.on("data", (data) => {
	console.log(`stdout: ${data}`);
});

action.on("close", (code) => {
	console.log(`child process close all stdio with code ${code}`);
});

action.on("exit", (code) => {
	console.log(`child process exited with code ${code}`);
});
