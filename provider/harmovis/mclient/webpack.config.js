

const resolve = require('path').resolve;
const webpack = require('webpack');
const HardSourceWebpackPlugin = require('hard-source-webpack-plugin');

module.exports = {
    output: {
	path: resolve(__dirname,'build'),
	filename: 'bundle.js'
    },
    devtool: 'source-map',
    module: {
	rules: [
	    {
		test: /\.js$/,
		loader: 'babel-loader',
		include: [resolve(__dirname, './src')],
		query: { "presets": ["@babel/react"] }
	    },
	    {
		test: /\.scss$/,
		use: ["style-loader", "css-loader", "sass-loader"]
	    },
	    {
		test: /\.css$/,
		use: [
		    'css-loader', // translates CSS into CommonJS
		    'sass-loader' // compiles Sass to CSS, using Node Sass by default
		]
	    }
	]
    },
    plugins: [
	// for compiling cache(speed up)
	new HardSourceWebpackPlugin(),
	// Optional: Enables reading mapbox token from environment variable
	new webpack.EnvironmentPlugin(['MAPBOX_ACCESS_TOKEN'])
    ]
};
