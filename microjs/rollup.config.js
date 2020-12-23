import babel from 'rollup-plugin-babel';
import babelrc from 'babelrc-rollup';
import html from 'rollup-plugin-html';

export default {
  entry: 'index.js',
  dest: 'micro.js',
  plugins: [
    html({
      include: '**/*.html'
    }),
    babel(babelrc()),
  ]
};
