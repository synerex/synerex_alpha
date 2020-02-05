
/*export interface Provider {
	Type: ProviderType
	Name: string
	Options: Option[]
	Running: boolean
}*/

export interface Command {
	Type: CommandType
	Name: string
	Option: Option["SET_AGENTS"] | Option["SET_AREA"] | Option["SET_CLOCK"] | Option["START_CLOCK"] | Option["STOP_CLOCK"]
}

export interface DialogStatus {
	Open: boolean; // 1サイクルにかかった時間
	Command?: Command;
}

export interface Log {
	ID: number
	Description: string
}


export interface RunningInfo {
	Duration: number,		// 1サイクルにかかった時間
	AgentsNum: number,
}

export interface Option {
	SET_AREA: {
		AreaCoord: Coord[]
	},
	SET_CLOCK: {
		Time: string
	},
	SET_AGENTS: {
		Type: AgentType,
		Num: number
	},
	START_CLOCK: {},
	STOP_CLOCK: {}
}

export interface Coord {
	Latitude: number,		// 1サイクルにかかった時間
	Longitude: number,
}

export enum AgentType {
	PEDESTRIAN,
	CAR,
}

export enum ContentType {
	LOG_MONITOR,
	AREA_MONITOR,
}

export enum CommandType {
	SET_AREA,
	SET_CLOCK,
	SET_AGENTS,
	START_CLOCK,
	STOP_CLOCK
}

export enum ProviderType {
	SCENARIO,
	VISUALIZATION,
	CAR,
	CLOCK,
	PEDESTRIAN,
	NODEID_SERVER,
	MONITOR_SERVER,
	SYNEREX_SERVER
}

export interface Order {

}

export class Provider {
	ID: number
	Type: ProviderType
	Logs: Log[]
	RunningInfos: RunningInfo[]

	constructor(id: number, type: ProviderType) {
		this.ID = id
		this.Type = type
		this.Logs = []
		this.RunningInfos = []
	}

	getName() {
		switch (this.Type) {
			case ProviderType.CAR:
				return "Car"
			case ProviderType.CLOCK:
				return "Clock"
			case ProviderType.SCENARIO:
				return "Scenario"
			case ProviderType.PEDESTRIAN:
				return "Pedestrian"
			case ProviderType.VISUALIZATION:
				return "Visualization"
			case ProviderType.NODEID_SERVER:
				return "NodeId Server"
			case ProviderType.MONITOR_SERVER:
				return "Monitor Server"
			case ProviderType.SYNEREX_SERVER:
				return "Synerex Server"
		}
	}

	getLogs() {
		return this.Logs
	}

	addLog(log: Log) {
		this.Logs.push(log)
	}

	addRunningInfo(rf: RunningInfo) {
		this.RunningInfos.push(rf)
	}

}