<template lang="pug">
  div
    v-data-table(
    v-model="selected",
    :items="snmpRules",
    :loading="pending",
    :headers="headers",
    item-key="_id",
    hide-actions,
    select-all
    )
      template(slot="items", slot-scope="props")
        tr
          td(@click.stop="")
            v-checkbox(v-model="props.selected", primary, hide-details)
          td
            snmp-rules-list-item-cell(
            :fields="['mibName', 'moduleName']",
            :source="props.item.oid"
            )
          td
            snmp-rules-list-item-cell(
            :fields="['value', 'regex']",
            :source="props.item.output"
            )
          td
            snmp-rules-list-item-cell(
            :fields="['value', 'regex']",
            :source="props.item.resource"
            )
          td
            snmp-rules-list-item-cell(
            :fields="['value', 'regex']",
            :source="props.item.component"
            )
          td
            snmp-rules-list-item-cell(
            v-if="props.item.state",
            :fields="['state', 'type']",
            :source="props.item.state"
            )
          td
            v-layout
              v-flex
                v-btn(icon, small)
                  v-icon edit
                v-btn.error--text(icon, small)
                  v-icon(color="error") delete
</template>

<script>
import SnmpRulesListItemCell from './snmp-rules-list-item-cell.vue';

export default {
  components: {
    SnmpRulesListItemCell,
  },
  props: {
    snmpRules: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: true,
    },
    showEditSnmpRuleModal: {
      type: Function,
      default: () => () => {},
    },
    showDeleteSnmpRuleModal: {
      type: Function,
      default: () => () => {},
    },
    showDeleteSelectedSnmpRuleModal: {
      type: Function,
      default: () => () => {},
    },
  },
  data() {
    return {
      selected: [],
    };
  },
  computed: {
    headers() {
      return [
        { text: 'oid', sortable: false },
        { text: 'output', sortable: false },
        { text: 'resource', sortable: false },
        { text: 'component', sortable: false },
        { text: 'state', sortable: false },
        { text: 'actions', sortable: false },
      ];
    },
  },
};
</script>
