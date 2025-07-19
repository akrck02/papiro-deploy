const { spawnSync } = require("child_process");

function chooseBinary() {
	return `main-linux-amd64`; // TODO: multi architecture
}

const binary = chooseBinary();
const mainScript = `${__dirname}/${binary}`;
const spawnSyncReturns = spawnSync(mainScript, {
	stdio: "inherit",
});
