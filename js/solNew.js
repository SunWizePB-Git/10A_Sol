import { WebSocketController, ViewMode } from './common.js'

export function run(prefix, url, viewMode) {
	const solNew = new solNew(prefix, url, viewMode)
}

class SolNew extends WebSocketController {

	open() {
		super.open()

		if (this.state.DeployParams === "") {
			return
		}

		this.show()
	}

	show() {
		this.showStatus()
		this.showSystem()
		this.showBattery()
		this.showLoadInfo()
		this.showSolar()
	}

	showStatus() {
		switch (this.viewMode) {
		case ViewMode.ViewFull:
			var status = document.getElementById("status")
			status.value = ""
			status.value += "Status:                      " + this.state.Status
			break;
		}
	}

	showSystem() {
		switch (this.viewMode) {
		case ViewMode.ViewFull:
			var ta = document.getElementById("system")
			ta.value = ""
			ta.value += "Configuration Mode:                      " + this.state.System.Config
			break;
		}
	}

	showBattery() {
		switch (this.viewMode) {
		case ViewMode.ViewFull:
			var ta = document.getElementById("battery")
			ta.value = ""
			ta.value += "* Percentage SOC:              " + this.state.Battery.Remaining + "\r\n"
			ta.value += "* Voltage (V):               " + this.state.Battery.Volts + "\r\n"
			ta.value += "* Current (A):               " + this.state.Battery.Amps
			break;
		case ViewMode.ViewTile:
			document.getElementById("battery-volts").innerText = this.state.Battery.Volts.toFixed(2)
			document.getElementById("battery-amps").innerText = this.state.Battery.Amps.toFixed(2)
			break;
		}
	}

	showLoadInfo() {
		switch (this.viewMode) {
		case ViewMode.ViewFull:
			var ta = document.getElementById("load")
			ta.value = ""
			ta.value += "* Voltage (V):               " + this.state.LoadInfo.Volts + "\r\n"
			ta.value += "* Current (A):               " + this.state.LoadInfo.Amps
			break;
		case ViewMode.ViewTile:
			document.getElementById("load-volts").innerText = this.state.LoadInfo.Volts.toFixed(2)
			document.getElementById("load-amps").innerText = this.state.LoadInfo.Amps.toFixed(2)
			break;
		}
	}

	showSolar() {
		switch (this.viewMode) {
		case ViewMode.ViewFull:
			var ta = document.getElementById("solar")
			ta.value = ""
			ta.value += "* Voltage (V):               " + this.state.Solar.Volts + "\r\n"
			ta.value += "* Current (A):               " + this.state.Solar.Amps
			break;
		case ViewMode.ViewTile:
			document.getElementById("solar-volts").innerText = this.state.Solar.Volts.toFixed(2)
			document.getElementById("solar-amps").innerText = this.state.Solar.Amps.toFixed(2)
			break;
		}
	}

	handle(msg) {
		switch(msg.Path) {
		case "update/status":
			this.state.Status = msg.Status
			this.showStatus()
			break
		case "update/system":
			this.state.System = msg.System
			this.showSystem()
			break
		case "update/battery":
			this.state.Battery = msg.Battery
			this.showBattery()
			break
		case "update/load":
			this.state.LoadInfo = msg.LoadInfo
			this.showLoadInfo()
			break
		case "update/solar":
			this.state.Solar = msg.Solar
			this.showSolar()
			break
		}
	}
}