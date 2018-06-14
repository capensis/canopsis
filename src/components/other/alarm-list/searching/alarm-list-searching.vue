<template lang="pug">
  v-toolbar.toolbar(dense, flat)
    v-text-field(
      label="Search"
      v-model="searchingText"
      hide-details
      single-line
      @keyup.enter="submit"
      @keyup.delete="clear"
    )
    v-btn(icon @click="submit")
      v-icon search
    v-btn(icon @click="clear")
      v-icon clear
    v-tooltip(bottom)
      v-btn(icon slot="activator")
        v-icon help_outline
      span Help on the advanced research :
      p  - [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt;
        |  [ AND|OR [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt; ]
      p The "-" before the research is required
      p Operators : <=, <, =, !=, >=, >, LIKE (For MongoDB regular expression)
      p Value's type : String between quote, Boolean ("TRUE", "FALSE"), Integer, Float, "NULL"
      dl
        dt Examples :
        dt - Connector = "connector_1"
        dd Alarms whose connectors are "connector_1"
        dt - Connector="connector_1" AND Resource="resource_3"
        dd Alarms whose connectors is "connector_1" and the ressources is "resource_3"
        dt - Connector="connector_1" OR Resource="resource_3"
        dd Alarms whose connectors is "connector_1" or the ressources is "resource_3"
        dt - Connector LIKE 1 OR Connector LIKE 2
        dd Alarms whose connectors contains 1 or 2
        dt - NOT Connector = "connector_1"
        dd Alarms whose connectors isn't "connector_1"

</template>

<script>
import omit from 'lodash/omit';

export default {
  name: 'alarm-list-searching',
  data() {
    return {
      searchingText: '',
    };
  },
  methods: {
    clear() {
      const query = omit(this.$route.query, ['search']);
      this.$router.push({
        query: {
          ...query,
        },

      });
    },
    submit() {
      const search = this.searchingText;
      this.$router.push({
        query: {
          ...this.$route.query,
          search,
        },

      });
    },
  },
};
</script>

<style scoped>
  .toolbar {
    background-color: white;
  }
</style>
