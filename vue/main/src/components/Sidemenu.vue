<template>
    <div>
        <el-menu name="side" :default-active="$route.path" @select="SelectHandler" router :unique-opened=false>
            <v-menu v-bind:items="this.itemstree"></v-menu>
        </el-menu>
    </div>
</template>

<script>
import { mapState } from 'vuex'
import VMenu from "./Menu";
    export default {
        name: "v-sidemenu",
        components: {
            VMenu
        },
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
                    else if (second == "route")  {
                        this.$store.commit("GetItemsTree");
                        this.$store.commit("GetAllItems");
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
            this.$store.commit("GetItemsTree");
            this.$store.commit("GetAllItems");
        },
        computed: {
            ...mapState([
                "itemstree",
                "envs"
            ])
        }
    };
</script>
