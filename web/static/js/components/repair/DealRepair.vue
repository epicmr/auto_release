<template>
    <div>
        <div>
            <el-input class="filter__input" v-on:input="OnInput" v-model="flt.user_id" placeholder="请输入用户ID" @keyup.enter.native="Query()" clearable></el-input>
            <el-button class="filter__button" type="primary" @click="Query()" accesskey="1">查询</el-button>
        </div>
        <div class="xx">
            <el-table :data="tableData" style="width: 100%">
                <el-table-column
                    type="selection"
                    width="55">
                </el-table-column>
                <el-table-column type="expand">
                    <template slot-scope="props">
                        <el-form label-position="left" inline>
                            <el-form-item label="订单 ID">
                                <span>{{ props.row.deal_id }}</span>
                            </el-form-item>
                            <el-form-item label="用户 ID">
                                <span>{{ props.row.user_id }}({{ props.row.user_phone }})</span>
                            </el-form-item>
                            <el-form-item label="患者名称">
                                <span>{{ props.row.patient_name }}</span>
                            </el-form-item>
                            <el-form-item label="门店 ID">
                                <span>{{ props.row.clinic_id }}</span>
                            </el-form-item>
                            <el-form-item label="门店名称">
                                <span>{{ props.row.clinic_name }}</span>
                            </el-form-item>
                        </el-form>
                    </template>
                </el-table-column>
                <div v-for="l in this.listCol">
                    <el-table-column :label="l.label" :width="l.width" :prop="l.prop" show-overflow-tooltip>
                    </el-table-column>
                </div>
            </el-table>
            <el-pagination
                background
                layout="prev, pager, next"
                @current-change="handleCurrentChange"
                :page-size="pageSize"
                :current-page.sync="curPage"
                :total="totalPage">
            </el-pagination>
        </div>
    </div>
</template>

<script>
export default {
    name:"v-dealrepair",
    data() {
        return {
            totalPage:1,
            curPage:1,
            pageSize:15,
            tableData: [],
            flt:{"user_id":"","page_size":"15","page_no":"1"},
            listCol: [
                {"label":"订单 ID","width":"100","prop":"deal_id"},
                {"label":"门店名称","width":"120","prop":"clinic_name"},
                {"label":"用户 ID","width":"100","prop":"user_id"},
                {"label":"患者姓名","width":"100","prop":"patient_name"},
                {"label":"预约日期","width":"100","prop":"reservation_start_time"},
                {"label":"医生姓名","width":"100","prop":"doctor_name"},
                {"label":"订单状态","width":"100","prop":"deal_state"},
                {"label":"应收金额","width":"100","prop":"receivable_price"},
                {"label":"实收金额","width":"100","prop":"price"},
                {"label":"支付状态","width":"100","prop":"pay_return_time"},
                {"label":"支付方式","width":"100","prop":"pay_type"},
                {"label":"订单来源","width":"100","prop":"source"}
            ]
        }
    },
    methods: {
        handleCurrentChange(val) {
            this.curPage = val
            this.Query()
        },
        async Query() {
            console.log(this.flt)
            await this.$http
                .post('http://cgi.gstyun.cn/cgi-bin/deal/temporaryList', this.flt)
                .then(response => {
                    if (parseInt(response.data.status, 10) == 0) {
                        console.log(this.totalPage)
                        this.totalPage = parseInt(response.data.total_page, 10)
                        this.tableData = response.data.id_list
                        console.log(this.tableData)
                    }
                })
        }
    }
}
</script>

<style lang="scss">
.filter 
    margin: 5px
    border-radius: 4px
    font-size: 14px
    &__input
        width: 200px
    &__button
        width: 100px
.xx
    .el-form-item
        margin-right: 0;
        margin-bottom: 0;
        width: 50%;
    label
        width: 90px
        color: #99a9bf
</style>
