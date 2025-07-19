function chooseBinary() {
	if (platform === "linux" && arch === "x64") {
		return `main-linux-amd64`;
	}
}

const binary = chooseBinary();
const mainScript = `${__dirname}/${binary}`;
const spawnSyncReturns = childProcess.spawnSync(mainScript, {
	stdio: "inherit",
});
