export class root extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
                id : "",
                idInput : "",
                errorMessage : "",
                errorTimeout : null,
                WS_CONNECTION: new WebSocket("ws://localhost:8443/ws"),
                constructMessage: (topic, payload) => {
                    return JSON.stringify({
                        topic: topic,
                        payload: payload,
                    });
                },
                setStateRoot: (state) => {
                    this.setState(state)
                },
                setErrorMessage: (message) => {
                    clearTimeout(this.state.errorTimeout)
                    this.setState({
                        errorMessage : message,
                        errorTimeout : setTimeout(() => {
                            this.setState({
                                errorMessage : "",
                            })
                        }, 5000)
                    })
                }
            },
            (this.state.WS_CONNECTION.onmessage = (event) => {
                let message = JSON.parse(event.data);
                switch (message.topic) {
                    case "connected":
                        this.state.setStateRoot({
                            id: message.payload,
                        });
                        break;
                    case "error":
                        let errorMessage = message.payload.split("->").reverse()[0]
                        this.state.setErrorMessage(errorMessage);
                        break;
                    default:
                        console.log("Unknown message topic: " + event.data);
                        break;
                }
            });
        this.state.WS_CONNECTION.onclose = () => {
            setTimeout(() => {
                if (this.state.WS_CONNECTION.readyState === WebSocket.CLOSED) {}
                window.location.reload();
            }, 2000);
        };
        this.state.WS_CONNECTION.onopen = () => {
            let myLoop = () => {
                this.state.WS_CONNECTION.send(this.state.constructMessage("heartbeat", ""));
                setTimeout(myLoop, 15 * 1000);
            };
            setTimeout(myLoop, 15 * 1000);
        };
    }

    render() {
        return React.createElement(
            "div", {
                id: "root",
                onContextMenu: (e) => {
                    e.preventDefault();
                },
                style: {
                    fontFamily: "sans-serif",
                    display: "flex",
                    flexDirection: "column",
                    justifyContent: "center",
                    alignItems: "center",
                },
            },
            React.createElement("div", {
                style: {
                    display: "flex",
                    flexDirection: "column",
                    justifyContent: "center",
                    alignItems: "center",
                },
            }, this.state.errorMessage || "\u00a0"),
            "your id: " + this.state.id,
            React.createElement("div", {
                    style: {
                        display: "flex",
                        flexDirection: "column",
                        justifyContent: "center",
                        alignItems: "center",
                    },
                },
                "enter another id or share your id"
            ),
            React.createElement("div", {
                    style: {
                        display: "flex",
                        flexDirection: "row",
                        justifyContent: "center",
                        alignItems: "center",
                    },
                },
                React.createElement("input", {
                    type: "text",
                    id: "input",
                    value: this.state.idInput,
                    onChange: (e) => {
                        this.state.setStateRoot({
                            idInput: e.target.value,
                        });
                    },
                    style: {
                        width: "100px",
                        height: "20px",
                    },
                }),
                React.createElement("button", {
                    onClick: () => {
                        this.state.WS_CONNECTION.send(
                           this.state.constructMessage("startGame", this.state.idInput)
                        );
                    },
                    style: {
                        width: "100px",
                        height: "20px",
                    },
                }, "start game")
            ),
        );
    }
}