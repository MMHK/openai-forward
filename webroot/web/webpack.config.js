// webpack.config.js
const path = require('path');
const {ProgressPlugin} = require('webpack');
const { VueLoaderPlugin } = require('vue-loader');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");
const HTMLInlineCSSWebpackPlugin = require('html-inline-css-webpack-plugin').default;
const HtmlInlineScriptPlugin = require('html-inline-script-webpack-plugin');
const TerserPlugin = require('terser-webpack-plugin');

IS_DEV_SERVER = process.env.NODE_ENV === 'development';

module.exports = {
  mode: IS_DEV_SERVER ? 'development' : 'production',
  entry: path.resolve(__dirname, "src/main.js"), // 入口 SCSS 文件
  output: {
    path: path.resolve(__dirname, '../ui'),
    clean: true,
    filename: IS_DEV_SERVER ? 'main.js' : '[name].min.js',
  },
  resolve: {
    extensions: ['.js', '.vue', '.json'],
  },
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: 'vue-loader',
      },
      {
        test: /\.scss$/,
        use: [
          IS_DEV_SERVER ? 'style-loader' : MiniCssExtractPlugin.loader, // 将 CSS 注入 DOM
          'css-loader',   // 转换 CSS 为 CommonJS
          {
            loader: 'sass-loader',
            options: {
              sassOptions: {
                silenceDeprecations: [
                  'legacy-js-api',
                  'import',
                  'function-units',
                  'slash-div',
                  'global-builtin'
                ],
              }
            }
          },
        ],
      },
      {
        test: /\.html$/,
        use: [
          'html-loader',
        ],
      },
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: [
          'babel-loader',
        ],
      },
      {
        test: /\.(woff|woff2|eot|ttf|otf|svg)$/,
        type: "asset",
        parser: {
          dataUrlCondition: {
            maxSize: 300 * 1024, // 300kb
          }
        },
        generator: {
          filename: 'assets/fonts/[hash][ext]'
        },
      },
      {
        test: /\.(png|jpe?g|gif)$/,
        type: "asset",
        parser: {
          dataUrlCondition: {
            maxSize: 300 * 1024, // 300kb
          }
        },
        generator: {
          filename: 'assets/img/[hash][ext]'
        }
      },
      {
        test: /\.css$/i,
        use: [
          IS_DEV_SERVER ? 'style-loader' : MiniCssExtractPlugin.loader, // 将 CSS 注入 DOM
          'css-loader',
          'postcss-loader', // 使用 PostCSS 添加自动前缀
        ],
      }
    ],
  },
  plugins: [
    ...(IS_DEV_SERVER ? [] : [
      new MiniCssExtractPlugin({
        filename: IS_DEV_SERVER ? '[name].css' : '[name].min.css',
      }),

      // new HTMLInlineCSSWebpackPlugin(),
      // new HtmlInlineScriptPlugin(),
    ]),
    new VueLoaderPlugin(),
    new ProgressPlugin(),
  ].concat([
    path.resolve(__dirname, 'public/index.html'),
  ].map(filePath => new HtmlWebpackPlugin({
      template: filePath,
      filename: path.basename(filePath),
      inject: "body",
      chunks: ['main'],
      minify: false,
  }))),
  devtool: IS_DEV_SERVER ? 'inline-source-map' : false,
  optimization: {
    minimizer: [
      new CssMinimizerPlugin({
        minimizerOptions: {
          preset: [
            'default',
            {
              discardComments: { removeAll: true },
              // 避免過度優化移除 Tailwind 類別
              discardUnused: false, // 禁用移除未使用的類別（如果需要）
            },
          ],
        },
      }),
      // 添加 JS 压缩插件
      new TerserPlugin({
        terserOptions: {
          compress: {
            drop_console: true, // 可选：移除console.log
          },
        },
      }),
    ],
    minimize: !IS_DEV_SERVER,
  },
  devServer: {
    static: {
      directory: path.join(__dirname, 'public'),
    },
    compress: true,
    port: 9000,
    hot: true,
    proxy: [
      {
        context: ['/api'],
        target: 'http://127.0.0.1:3005',
        changeOrigin: true,
        secure: false,
        headers: {
          Connection: 'keep-alive',
        },
      },
    ]
  },
};
