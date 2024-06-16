
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
            "game"
		)
    }
}