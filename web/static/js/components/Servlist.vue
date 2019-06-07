<template>
    <div>
        <el-row>
            <el-checkbox :indeterminate="isIndeterminate" v-model="checkAll" @change="handleCheckAllChange">全选</el-checkbox>
        </el-row>
        <el-checkbox-group v-model="checkedApp" v-if="list_type === parseInt(serv.serv_type)" v-for="serv in this.servs">
            <el-checkbox name="serv_list" :label="serv.serv_name" :key="serv.serv_name" style="margin-top: 10px" size="medium"></el-checkbox>
        </el-checkbox-group>
    </div>
</template>

<script>
    import { mapState } from 'vuex'

    export default {
        name: "v-servlist",
        props : {
            list_type : {
                type : Number,
                default : 0
            }
        },
        data() {
            return {
                parent:this,
                checkAll:false,
                isIndeterminate:false,
                checkedApp:[]
            }
        },
        methods: {
            handleCheckAllChange(val) {
                this.checkedApp = []
                if (val) {
                    for (let i in this.servs) {
                        this.checkedApp = this.checkedApp.concat(this.servs[i].serv_name)
                    }
                }
                this.isIndeterminate = false;
            },
        },
        computed: {
            ...mapState([
                "servs"
            ])
        }
    }
</script>
<style scoped>
small{
    font-size: 12px;
    color: #66b1ff;
    vertical-align: bottom;
    display: inline-block;
    margin-top: 10px;
}
</style>

