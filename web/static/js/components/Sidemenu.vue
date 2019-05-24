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
                <el-menu-item v-for="env in this.envs" :key="env.id" :index="'release-' + env.name">
                    <i class="el-icon-notebook-1"></i>
                    <span slot="title">{{env.name}}</span>
                </el-menu-item>
            </el-submenu>

            <el-submenu index="conf">
                <template slot="title">
                    <i class="el-icon-setting"></i>
                    <span slot="title">配置</span>
                </template>
                <el-menu-item index="conf-env">
                    <i class="el-icon-picture"></i>
                    <span slot="title">环境</span>
                </el-menu-item>
                <el-menu-item index="conf-host">
                    <i class="el-icon-monitor"></i>
                    <span slot="title">主机</span>
                </el-menu-item>
                <el-menu-item index="conf">
                    <i class="el-icon-s-grid"></i>
                    <span slot="title">服务</span>
                </el-menu-item>
            </el-submenu>

            <el-submenu index="repair">
                <template slot="title">
                    <i class="el-icon-s-open"></i>
                    <span slot="title">修复</span>
                </template>
                <el-menu-item index="repair-user">
                    <i class="el-icon-s-help"></i>
                    <span slot="title">删除用户</span>
                </el-menu-item>
            </el-submenu>

        </el-menu>
    </div>
</template>

<script>
    import { mapState } from 'vuex'

    export default {
        name: "v-sidemenu",
        data() {
            return {
            }
        },
        methods: {
            SelectHandler(index, indexpath) {
                let first = index.split('-')[0]
                let second = index.split('-')[1]
                if (first == 'conf') {
                    if (second == "env")  {
                        this.$store.commit("GetEnvs", {"info_type" : 1})
                    }
                    else if (second == "host")  {
                        this.$store.commit("GetHosts");
                    }
                    else {
                        this.$store.commit("GetConfs");
                    }
                }
                else if (first === 'release') {
                    this.$store.commit("GetServs", {"env":second});
                }
            }
        },
        mounted() {
            this.$store.commit("GetEnvs")
            this.$store.commit("GetHosts");
            this.$store.commit("GetConfs");
        },
        computed: {
            ...mapState([
                "envs"
            ])
        }
    };
</script>
