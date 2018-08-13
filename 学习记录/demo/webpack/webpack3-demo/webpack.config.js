const path = require('path')
const UglifyJSPlugin = require('uglifyjs-webpack-plugin')

module.exports = {
  entry: './src2/index.js',
  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'dist')
  },
  plugins: [
    // new UglifyJSPlugin()
    // new UglifyJSPlugin({
    //   uglifyOptions: {
    //     compress: {
    //       pure_funcs: ['Math.floor']
    //     }
    //   }
    // })
  ],
  devtool: false
}
