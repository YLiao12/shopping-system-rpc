<template>
  <div>
    <div class="user">
    <h6 class="userinfo" v :title="user_info">{{user_info}}</h6>
    </div>
    <div class = "search">
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
    
  </div>
  
</template>


<script>
import axios from 'axios'


export default {
  name: "ShoppingProducts",
  data() {
    return {
      result:"",
      user_info:"",
      products:[],
      multipleSelection:[],
      input_pid:"",
    }
  },
  mounted() {
    this.AllProduct()
    this.GetUser()
  },
  methods: {
    AllProduct() {
        axios.get("http://103.49.160.227:5555/AllProducts")
        .then((response)=>{
          if(response.status == 200){
            this.products = response.data;
            console.log(response.data)
            
          }
        })
      
    },
    GetUser(){
      axios.get("http://103.49.160.227:5555/GetUser")
      .then((response)=>{
          if(response.status == 200){
            var user = response.data;
            console.log(response.data)
            this.user_info = "Hi,"+user.name+"  your balance now is: "+user.balance
            
          }
        })
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
      var u_id = this.user_id
      axios.get("http://103.49.160.227:5555/MakeOrder?productId="+p_id+"&userId="+u_id)
      .then((response)=>{
          if(response.status == 200){
            this.AllProduct()
          }
        })
    },
    OneProduct(){
      var input_product = this.input_pid
      axios.get("http://103.49.160.227:5555/OneProduct?productId="+input_product)
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

</style>