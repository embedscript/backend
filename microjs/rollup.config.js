import babel from 'rollup-plugin-babel';
import babelrc from 'babelrc-rollup';

export default {
  entry: 'index.js',
  dest: 'micro.js',
  plugins: [
    babel(babelrc())
  ]
};
