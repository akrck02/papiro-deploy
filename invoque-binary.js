function chooseBinary() {
	return `main-linux-amd64`; // TODO: multi architecture
}

const binary = chooseBinary();
const mainScript = `${__dirname}/${binary}`;
const spawnSyncReturns = childProcess.spawnSync(mainScript, {
	stdio: "inherit",
});
