const { spawnSync } = require("child_process");

function chooseBinary() {
	return `main-linux-amd64`;
}

const binary = chooseBinary();
const mainScript = `${__dirname}/${binary}`;
const action = spawnSync(mainScript, []);
console.log(action.stdout.toString());
