<template>
    <div>
        <input type="text" class="form-control" placeholder="Search" v-model="search">
        <paginate
            name="paginatedAlarms"
            :list="filteredAlarms"
            :per="3"
            tag="div"
        >
            <table>
                <tr>
                    <th v-on:click="sortColumn('component')">component</th>
                    <th v-on:click="sortColumn('resource')">resource</th>
                    <th>creation_date</th>
                    <th>state</th>
                    <th>status</th>
                    <th>actions</th>
                </tr>

                    <tr v-for="item in paginated('paginatedAlarms')">
                        <td>{{ item.v.component }}</td>
                        <td>{{ item.v.resource }}</td>
                        <td>{{ item.v.creation_date | ts2date}}</td>
                        <td v-bind:class="state2color(item.v.state.val)">{{ item.v.state.val }}</td>
                        <td v-bind:class="state2color(item.v.status.val)">{{ item.v.status.val }}</td>
                        <td><button v-on:click="deleteRecord(item._id)">X</button></td>
                    </tr>

            </table>
        </paginate>
        <paginate-links
          for="paginatedAlarms"
          :show-step-links="true"
        ></paginate-links>
    </div>
</template>

<script>

import { ALARMS } from './dataset.js';

export default {
    data () {
        return {
            alarms: ALARMS,
            filteredAlarms: ALARMS,
            paginate: ['paginatedAlarms'],
            sort: {name: "", order: "ASC"},
            search: ""
        }
    },
    methods:{
        state2color: function(value){
            switch(value){
                case 0:
                    return "green";
                case 1:
                    return "yellow";
                case 2:
                    return "orange";
                case 3:
                    return "red";
                default:
                    return "grey";
            }
        },
        sortColumn: function(column){

            if(this.sort.name === column && this.sort.order === "ASC"){
                this.filteredAlarms.sort(
                    (a, b) => a.v[column].localeCompare(b.v[column])
                ).reverse()
                this.sort.order = "DESC"
            }else{
                this.filteredAlarms.sort(
                    (a, b) => a.v[column].localeCompare(b.v[column])
                )
                this.sort.order = "ASC"
            }

            this.sort.name = column
        },
        deleteRecord(truc){
            console.error(truc)
        }
    },
    filters:{
        ts2date: function(value){
            let d = new Date(value*1000)
            return d.getDate() + '/' + (d.getMonth()+1) + '/' + d.getFullYear();
        }
    },
    watch: {
        search: function() {
            let me = this
            me.filteredAlarms = me.alarms.filter(function(item) {
                console.error(item)
                let searchRegex = new RegExp(me.search, 'i')
                if(searchRegex.test(item.v.component))
                        return true
                if(searchRegex.test(item.v.resource))
                    return true;
                return false
            })
        }
    }
}
</script>

<style  scoped>
    table{
        width: 100%;
    }

    .green { color:green; }

    .yellow { color:yellow; }

    .orange { color:orange; }

    .red { color:red; }

    .grey { color:grey; }
</style>
