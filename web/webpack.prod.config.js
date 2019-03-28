const os = require('os');
const path = require('path');
const webpack = require('webpack')
const pkg = require('./package.json')
const projectRoot = path.resolve(__dirname);
const cleanDistFolderPlugin = require('clean-webpack-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const ExtractTextPlugin = require('extract-text-webpack-plugin')
const notifier = require('node-notifier')
const WriteFileWebPackPlugin = require('write-file-webpack-plugin')
const FriendlyErrorsPlugin = require('friendly-errors-webpack-plugin')
const UglifyJSPlugin = require('uglifyjs-webpack-plugin')
const CompressionPlugin = require("compression-webpack-plugin");
const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin;
const ParallelUglifyPlugin = require('webpack-parallel-uglify-plugin');

const devServerHost = process.env.HOST || '120.25.154.225'
const devServerPort = process.env.PORT || '8000'

module.exports = {
    entry: __dirname + "/static/js/index.js",
    output: {
            path: __dirname + "/dist",
            filename: "bundle.js"
    },
    devtool: 'source-map',
    resolve: {
        alias: {
            'vue$': 'vue/dist/vue.esm.js',
            'src': path.resolve(__dirname, 'static/js'),
            'css': path.resolve(__dirname, 'static/css/index.scss'),
            'components': path.resolve(__dirname, 'static/js/components')
        },
        extensions: ['*', '.js', '.vue', '.json']
    },
    devServer: {
        host: devServerHost,
        port: devServerPort,
        contentBase: path.join(__dirname, "dist"),
    },
    module: {
        rules: [
            {
                test: /\.vue$/,
                loader: 'vue-loader',
                options: {
                    loaders: {
                        css: ExtractTextPlugin.extract({
                            use: [{
                                loader: 'css-loader',
                                options: {
                                    minimize: true,
                                    sourceMap: false
                                }
                            }],
                            fallback: 'vue-style-loader'
                        }),
                        scss: ExtractTextPlugin.extract({
                            use: [{
                                loader: 'css-loader',
                                options: {
                                    minimize: true,
                                    sourceMap: false
                                }
                            }, {
                                loader: 'sass-loader',
                                options: {
                                    sourceMap: false,
                                    indentedSyntax: true
                                }
                            }],
                            fallback: 'vue-style-loader'
                        }),
                    }
                }
            },
            {
                test: /\.js$/,
                use: 'babel-loader',
                exclude: /node_modules/
            },
            {
                test: /\.(css|scss)$/,
                use: ExtractTextPlugin.extract({
                    use: [{
                        loader: 'css-loader',
                        options: {
                            minimize: true,
                            sourceMap: false
                        }
                    },
                    {
                        loader: 'sass-loader',
                        options: {
                            sourceMap: false
                        }
                    }],
                    fallback: 'vue-style-loader'
                })
            },
            {
                test: /\.(png|jpe?g|gif|svg)(\?.*)?$/,
                loader: 'url-loader',
                options: {
                    limit: 10000,
                    name: path.posix.join("dist", 'img/[name].bundle.[ext]')
                }
            },
            {
                test: /\.(mp4|webm|ogg|mp3|wav|flac|aac)(\?.*)?$/,
                loader: 'url-loader',
                options: {
                    limit: 10000,
                    name: path.posix.join("dist", 'media/[name].bundle.[ext]')
                }
            },
            {
                test: /\.(woff2?|eot|ttf|otf)(\?.*)?$/,
                loader: 'url-loader',
                options: {
                    limit: 10000,
                    name: path.posix.join("dist", 'fonts/[name].bundle.[ext]')
                }
            }
        ]
    },
    plugins: [
        new webpack.DefinePlugin({ 'process.env': { NODE_ENV: JSON.stringify(process.env.NODE_ENV) } }), // production | development
        new cleanDistFolderPlugin(['dist']),
        new webpack.NoEmitOnErrorsPlugin(),
        new ExtractTextPlugin({ filename: "[name].css" }),
        new HtmlWebpackPlugin({
            template: path.resolve(__dirname, 'static/index.html'),
            filename: 'index.html',
            inject: true,
            minify: {
                removeComments: false,
                collapseWhitespace: false,
                removeAttributeQuotes: false
            },
            chunks: ['main'],
            chunksSortMode: 'dependency',
            hash: true,
            showErrors: true
        }),
        new WriteFileWebPackPlugin(),
        new CompressionPlugin({
            asset: '[path].gz[query]', //目标资源名称。[file] 会被替换成原资源。[path] 会被替换成原资源路径，[query] 替换成原查询字符串
            algorithm: 'gzip',//算法
            test: /\.(js|html)$/,
            threshold: 10240,//只处理比这个值大的资源。按字节计算
            minRatio: 0.8//只有压缩率比这个值小的资源才会被处理
        }),
        new BundleAnalyzerPlugin(
            {
                analyzerMode: 'server',
                analyzerHost: '120.25.154.225',
                analyzerPort: 8889,
                reportFilename: 'report.html',
                defaultSizes: 'parsed',
                openAnalyzer: true,
                generateStatsFile: false,
                statsFilename: 'stats.json',
                statsOptions: null,
                logLevel: 'info'
            }
        )
    ]
};
