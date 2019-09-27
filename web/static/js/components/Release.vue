<template>
    <div>
        <form>
            <el-row>
                <el-col :span=8>
                    <v-servlist v-bind:list_type="1"></v-servlist>
                </el-col>
                <el-col :span=8>
                    <v-servlist v-bind:list_type="2"></v-servlist>
                </el-col>
                <el-col :span=8>
                    <v-servlist v-bind:list_type="3"></v-servlist>
                </el-col>
            </el-row>
            <el-row>
                <v-button></v-button>
            </el-row>
        </form>
    </div>
</template>

<script>
import VServlist from "./Servlist";
import VButton from "./Button";

export default {
    components: {
        VServlist,
        VButton
    },
    name: "v-release",
    data() {
        return {
            parent:this
        }
    }
}
</script>
