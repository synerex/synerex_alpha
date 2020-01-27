import { createMuiTheme } from '@material-ui/core/styles'

const theme = createMuiTheme({
	// Material-UIコンポーネントのclassのstyleを上書きする
	overrides: {
		MuiButton: {
			root: {
				// ボタン内アルファベット文字を大文字変換しない
				textTransform: 'none',
			},
		},
	},
})

export default theme