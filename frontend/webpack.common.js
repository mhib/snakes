const path = require('path');

module.exports = {
  entry: {
    game: './game.jsx',
    lobby: './lobby.jsx',
  },
  output: {
    filename: '[name].js',
    path: path.resolve(__dirname, 'public'),
  },
  optimization: {
    splitChunks: {
      cacheGroups: {
        commons: {
          name: 'commons',
          chunks: 'initial',
          minChunks: 2,
        },
      },
    },
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        loader: 'babel-loader',
      },
    ],
  },
  resolve: {
    extensions: ['*', '.js', '.jsx'],
  },
};
