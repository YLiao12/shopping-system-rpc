<template>
  <div>
    <div class="user">
    <h6 class="userinfo" v :title="user_info">{{user_info}}</h6>
    </div>
    <div class = "search">
      <el-input class="usersearch" v-model="user_id" placeholder="1～10">
      <template slot="prepend">UserId:</template>
      <el-button slot="append" icon="el-icon-search" @click="GetUser()"></el-button> 
      </el-input>
      <el-input  v-model="input_pid" placeholder="please enter productId">
        <el-button slot="append" icon="el-icon-search" @click="OneProduct()"></el-button> 
      </el-input>
    </div>
    <div class = "products-list" >
       <el-table
        :data="products"
        tooltip-effect="dark"
        style="width: 48%; margin:auto;"
        
        @selection-change="handleSelectionChange">
        <el-table-column
          type="selection"
          width="40">
        </el-table-column>
        <el-table-column
          prop="name"
          label="Product_Name"
          width="160"
          align="center">
        </el-table-column>
        <el-table-column
          prop="stock"
          label="Product_Stock"
          width="160"
          align="center">
        </el-table-column>
        <el-table-column
          prop="price"
          label="Product_Price"
          align="center">
        </el-table-column>
        <el-table-column
        width="80"
        align="right">
        <template slot="header">
          <el-button type="primary" size="small" @click="AllProduct()">All</el-button>
        </template>
        </el-table-column>
       <el-table-column
       width="80"
        align="right">
        <template slot="header">
          <el-button type="primary" size="small" @click="BuyProduct()">Buy</el-button>
        </template>
        </el-table-column>
      </el-table>
    </div>
    <div>
      <p class="hint">System now can just handle the order with one product</p>
    </div>
    
  </div>
  
</template>


<script>
import axios from 'axios'

var RPC_URL = "http://139.196.187.198:5555"

export default {
  name: "ShoppingProducts",
  data() {
    return {
      result:"",
      user_info:"",
      products:[],
      multipleSelection:[],
      input_pid:"",
      oreder_status:0,
      user_id:"1",
    }
  },
  mounted() {
    this.AllProduct()
    this.GetUser()
  },
  methods: {
    AllProduct() {
        axios.get(RPC_URL+"/AllProducts")
        .then((response)=>{
          if(response.status == 200){
            var data = response.data
            this.products = response.data;
            for(let i in data){
              if(!Object.prototype.hasOwnProperty.call(data[i],"stock")){
                Object.assign(data[i],{stock:0})
              }
            }
           // console.log(response.data)
            
          }
        })
      
    },
    GetUser(){
      var userid = this.user_id;
      var userint = parseInt(userid);
      if(userint<=10 && userint>=1){
        axios.get(RPC_URL+"/GetUser?userId="+userid)
        .then((response)=>{
          if(response.status == 200){
            var user = response.data;
            console.log(response.data)
            this.user_info = "Hi,"+user.name+"  your balance now is: "+user.balance
          }
        })
      }else{
        this.$message('can only input userId in 1~10');
      }
      
    },
    handleSelectionChange(val) {
        console.log(val)
        this.multipleSelection = val;
      },
    
    BuyProduct(){
      if(this.multipleSelection.length == 0){
        alert("you need to choose at least one product");
      }
      var p_id = this.multipleSelection[0].id
      var u_id = "1"
      axios.get(RPC_URL+"/MakeOrder?productId="+p_id+"&userId="+u_id)
      .then((response)=>{
          if(response.status == 200){
            var order_result = response.data;
            console.log(order_result)
            if(order_result == 2){
              this.order_result = 0
              this.$message('Failed to order,understock');
            }
            if(order_result == 1){
              this.order_result = 0
              this.$message('Failed to order, insufficient balance');
            }
            this.AllProduct()
            this.GetUser()
          }
        })
    },
    OneProduct(){
      var input_product = this.input_pid
      axios.get(RPC_URL+"/OneProduct?productId="+input_product)
      .then((response)=>{
          if(response.status == 200){
            this.products = response.data;
            console.log(response.data)
            
          }
        })
      this.input_pid=""
    }

       
  }
}
</script>


<style scoped>

.select{
  margin-top: 100px;
}
.search{
  width:40%;
  margin: auto;
}
.usersearch{
  width:50%
}

.products-list{
    text-align: center;
}
.el-row {
    margin: 0, auto;
    margin-bottom: 20px;
    text-align: center;
    &:last-child {
      margin-bottom: 0;
    }
  }
  .el-col {
    border-radius: 4px;
    text-align: center;
    
  }
  .bg-purple-dark {
    background: #99a9bf;
    text-align: center;
  }
  .bg-purple {
    background: #d3dce6;
  }
  .bg-purple-light {
    background: #e5e9f2;
  }
  .grid-content {
    min-height: 36px;
  }
  .row-bg {
    padding: 10px 0;
  }
  .hint{
    font-family: Arial;
    font-size:9px;
    color:#CDCDCD;
  }

</style>