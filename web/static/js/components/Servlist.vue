<template>
    <div>
        <el-checkbox-group v-model="checkedApp">
            <el-row v-if="list_type === parseInt(serv.serv_type)" v-for="serv in this.$store.state.data">
                <el-col :span=12>
                    <el-checkbox name="serv_list" :label="serv.serv_name" :key="serv.serv_name" style="margin-top: 10px" size="medium"></el-checkbox>
                </el-col>
                <el-col :span=12>
                    <el-row v-for="time in serv.serv_state">
                        <el-col :span=12>
                            <small>{{time.host_name}}</small>
                        </el-col>
                        <el-col :span=12>
                            <small>{{time.serv_time}}</small>
                        </el-col>
                    </el-row>
                </el-col>
            </el-row>
        </el-checkbox-group>
    </div>
</template>

<script>
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
                checkedApp:[]
            }
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

