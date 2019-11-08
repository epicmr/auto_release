<template>
    <div>
        <div class="filter">
            <el-input class="input" v-model="userId" placeholder="请输入用户ID" @keyup.enter.native="query()" clearable></el-input>
            <el-button class="button" type="primary" @click="query()" accesskey="1">查询</el-button>
        </div>
        <div class="table">
            <el-pagination v-if="this.pageSize > 10"
                @size-change="handleSizeChange"
                :page-sizes="[10, 15, 30, 50, 100]"
                :total="totalNum"
                layout="sizes, total">
            </el-pagination>
            <el-table :data="tableData" style="width: 100%">
                <el-table-column type="selection" width="50">
                </el-table-column>
                <el-table-column label="订单 ID" :width="100" prop="deal_id" type="text" :formatter="pData" show-overflow-tooltip>
                </el-table-column>
                <el-table-column label="订单类型" :width="76" prop="deal_type" type="deal_type" :formatter="pData" show-overflow-tooltip>
                </el-table-column>
                <el-table-column label="门店名称" :width="120" prop="clinic_name" type="text" :formatter="pData" show-overflow-tooltip>
                </el-table-column>
                <el-table-column label="用户 ID" :width="136" prop="user_id" type="text" :formatter="pData" show-overflow-tooltip>
                </el-table-column>
                <el-table-column label="患者姓名" :width="76" prop="patient_name" type="text" :formatter="pData" show-overflow-tooltip>
                </el-table-column>
                <el-table-column label="预约日期" :width="100" prop="reservation_start_time" type="unix2date" :formatter="pData" show-overflow-tooltip>
                </el-table-column>
                <el-table-column label="服务时间" :width="100" show-overflow-tooltip>
                    <template slot-scope="props">
                        <span>{{resvTime(props.row.reservation_start_time, props.row.reservation_end_time)}}</span>
                    </template>
                </el-table-column>
                <el-table-column label="医生姓名" :width="76" prop="doctor_name" type="text" :formatter="pData" show-overflow-tooltip>
                </el-table-column>
                <el-table-column label="订单状态" :width="76" prop="deal_state" type="text" :formatter="pData" show-overflow-tooltip>
                    <template slot-scope="props">
                        <el-tag type="success" class="tag" effect="dark">{{props.row.deal_state}}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column label="应收金额" :width="76" prop="receivable_price" type="price" :formatter="pData" show-overflow-tooltip>
                </el-table-column>
                <el-table-column label="实收金额" :width="76" prop="price" type="price" :formatter="pData" show-overflow-tooltip>
                </el-table-column>
                <el-table-column label="支付状态" :width="76" prop="pay_return_time" type="pay_status" :formatter="pData" show-overflow-tooltip>
                </el-table-column>
                <el-table-column label="支付方式" :width="76" prop="pay_type" type="pay_type" :formatter="pData" show-overflow-tooltip>
                </el-table-column>
                <el-table-column label="订单来源" :width="76" prop="source" type="source" :formatter="pData" show-overflow-tooltip>
                </el-table-column>
                <el-table-column type="expand" class="expand">
                    <template slot-scope="props">
                        <el-form label-position="left" inline class="demo-table-expand">
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
                <el-table-column label="操作">
                    <template slot-scope="scope">
                        <el-button size="mini" plain type="warning" v-if="checkOperation(scope.row) === 1"
                                               @click="submit(scope.row)">退 号</el-button>
                        <el-button size="mini" plain type="warning" v-if="checkOperation(scope.row) === 2"
                                               @click="submit(scope.row)">退 款</el-button>
                        <el-button size="mini" plain type="warning" v-if="checkOperation(scope.row) === 3"
                                               @click="submit(scope.row)">退 款</el-button>
                        <el-button size="mini" plain type="warning" v-if="checkOperation(scope.row) === 4"
                                               @click="submit(scope.row)">回 滚</el-button>
                    </template>
                </el-table-column>
            </el-table>
            <el-pagination
                background
                @current-change="handleCurrentChange"
                :page-size="this.pageSize"
                :current-page.sync="curPage"
                @size-change="handleSizeChange"
                :page-sizes="[10, 15, 30, 50, 100]"
                layout="sizes, prev, pager, next, total"
                :total="totalNum">
            </el-pagination>
        </div>
    </div>
</template>

<script>
export default {
    name:"v-dealrepair",
    data() {
        return {
            totalNum:1,
            userId:'',
            curPage:1,
            pageSize:10,
            tableData:[],
            dealType: {
                1:"处方订单",
                2:"拍照处方",
                3:"挂号订单",
                4:"预约订单",
                5:"排队订单",
                6:"活动订单",
                7:"商城订单",
                8:"咨询订单",
                9:"贵细订单",
                10:"会员订单",
                11:"补代煎单"
            }
        }
    },
    methods: {
        handleCurrentChange(val) {
            this.curPage = val
            this.query()
        },
        handleSizeChange(val) {
            this.pageSize = val
            this.query()
        },
        checkOperation(deal) {
            if (deal.cancel_time > 0) 
                return  0
            else if (deal.deal_type == 3 || deal.deal_type == 4)
                if (deal.arrived_time > 0)
                    return 1//退号
                else
                    return 3//预约退款
            else if (deal.deal_type == 1)
                if (deal.pay_return_time > 0)
                    return 2//处方退款
                else
                    return 4//处方回滚
            else
                return 0
        },
        resvTime(val1,val2){
            return new Date(val1 * 1000).format("hh:mm") + "-" + new Date(val2 * 1000).format("hh:mm");
        },
        pData(row, column, cellValue, index) {
            if (column.type == "unix2date") {
                return new Date(cellValue * 1000).format("yyyy-MM-dd");  
            }
            else if (column.type == "price") {
                return cellValue / 10000;  
            }
            else if (column.type == "unix2hour") {
                return new Date(cellValue * 1000).format("hh:mm");  
            }
            else if (column.type == "pay_status") {
                return cellValue > 0 ? "已付费" : "未付费";
            }
            else if (column.type == "deal_type") {
                return this.dealType[cellValue];
            }
            return cellValue
        },
        query() {
            this.$http
                .post('http://cgi.gstyun.cn/cgi-bin/deal/temporaryList', {"user_id":this.userId, "page_no":this.curPage+"", "page_size":this.pageSize+""})
                .then(response => {
                    if (parseInt(response.data.status, 10) == 0) {
                        this.totalNum = parseInt(response.data.total_num, 10)
                        this.tableData = response.data.id_list
                    }
                })
                .catch(err => {
                    console.info("%c [axiso catch error]", "color:orange", err)
                });
        },
        refund(deal) {
            let t = this.checkOperation(deal)
            if (t === 1) {
                console.log("退号")
            } else if (t === 2 || t === 3) {
                console.log("退款")
            } else if (t === 4)
                console.log("回滚")
            else{
            }
        },
        submit(deal) {
            this.$prompt('请输入工号', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                inputErrorMessage: '工号不正确'
            }).then(({ value }) => {
                this.refund(deal)
                this.$message({
                    type: 'success',
                    message: '操作人工号: ' + value
                });
            }).catch((err) => {
                console.log(err)
                this.$message({
                    type: 'info',
                    message: '放弃操作'
                });       
            });
        }
    }
}
</script>

<style lang="scss">
.filter {
    margin: 5px;
    border-radius: 4px;
    font-size: 14px;
    .input {
        width: 200px;
    }
    .button {
        width: 80px;
    }
}

.table {
    .demo-table-expand {
        font-size: 0;
    }

    .demo-table-expand label {
        width: 90px;
        color: #99a9bf;
    }

    .demo-table-expand .el-form-item {
        margin-right: 0;
        margin-bottom: 0;
        width: 30%;
    }
}

</style>
