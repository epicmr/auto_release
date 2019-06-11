<template>
    <div>
        <label v-for="item in this.items">
            <el-submenu :index="item.index" v-if="Array.isArray(item.items) && item.items.length > 0">
                <template slot="title">
                    <i :class="item.class"></i>
                    <span slot="title">{{item.name}}</span>
                </template>
                <v-menu v-bind:items="item.items"></v-menu>
            </el-submenu>
            <el-menu-item :index="item.index" v-else>
                <i :class="item.class"></i>
                <span slot="title">{{item.name}}</span>
            </el-menu-item>
        </label>
    </div>
</template>

<script>
import VMenu from "./Menu";

    export default {
        name: "v-menu",
        props : {
            items : {
                type : Array,
                default : []
            }
        },
        data() {
            return {
            }
        },
        methods: {
        },
        mounted() {
        },
        computed: {
        }
    };
</script>
