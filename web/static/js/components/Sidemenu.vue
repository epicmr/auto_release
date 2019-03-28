<template>
    <div>
        <h5></h5>
        <el-menu name="side" :default-active="$route.path" @select="SelectHandler" router class="el-menu-vertical-demo" :unique-opened=true>
            <el-menu-item index="menu" disabled>
                <i class="el-icon-menu"></i>
                <span slot="title">菜单</span>
            </el-menu-item>

            <el-submenu index="release">
                <template slot="title">
                    <i class="el-icon-upload"></i>
                    <span slot="title">发布</span>
                </template>
                <el-menu-item :index="'release-' + name" v-for="(host, name) in hosts">
                    <i class="el-icon-upload"></i>
                    <span slot="title">{{name}}</span>
                </el-menu-item>
            </el-submenu>

            <el-submenu index="conf">
                <template slot="title">
                    <i class="el-icon-conf"></i>
                    <span slot="title">配置</span>
                </template>
                <el-menu-item index="conf">
                    <i class="el-icon-conf"></i>
                    <span slot="title">服务</span>
                </el-menu-item>
            </el-submenu>

        </el-menu>
    </div>
</template>

<script>
    export default {
        name: "v-sidemenu",
        data() {
            return {
                hosts: []
            }
        },
        methods: {
            SelectHandler(index, indexpath) {
                let env = index.split('-')[0]
                if (env == 'conf') {
                    this.$store.dispatch("GetConfs");
                }
                else if (env === 'release') {
                    this.$store.dispatch("GetServs", {"env":index.split('-')[1]});
                }
                else {
                    console.log(this.hosts)
                }
            }
        },
        mounted() {
            let _this = this
            console.log("OnSideMenuCreate")
            this.$http
                .get('/api/hosts')
                .then(response => {
                    _this.hosts = response.data
                })
            console.log(_this.hosts)
        }
    };
</script>
