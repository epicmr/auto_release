<template>
    <div>
        <el-row>
            <el-col :span=4 class="site-header">
            </el-col>
            <el-col :span=16 id="header" class="site-header botton">
            </el-col div>
        </el-row>
    </div>
</template>

<script>
export default {
    name: "v-header",
    components: {
    },
    data() {
        return {
        }
    },
    methods: {
    }
};
</script>

<style>
.botton{
    border-bottom: 1px solid #dcdfe6;
}
.site-header{
    height: 5rem;
    margin: 0 auto;
}
</style>
