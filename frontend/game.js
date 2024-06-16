import { chessBoard } from "./chessBoard.js"

export class game extends React.Component {
	constructor(props) {
		super(props)
	}

    render() {
	
        return React.createElement("div", {
				style : {
					gap : "1vmin",
					position: "relative",
					marginTop : "1vmin",
					display: "flex",
					flexDirection : "column",
					alignItems : "center",
					justifyContent : "center",
				}
			},
			React.createElement(chessBoard, this.props),
			React.createElement("button", {
					style : {
						marginTop : "1vmin",
					},
					onClick : () => {
						this.props.WS_CONNECTION.send(this.props.constructMessage("endGame", ""))
					}
				}, "End Game",
			)
		)
    }
}