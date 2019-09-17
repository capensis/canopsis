<template lang="pug">
  div
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('heartbeat.title') }}
    .white
      v-layout(row, wrap)
        v-flex(v-show="hasDeleteAnyHeartbeatAccess && selected.length", xs4)
          v-btn(icon, @click="showDeleteSelectedHeartbeatsModal(selected)")
            v-icon(color="error") delete
      v-data-table(
        v-model="selected",
        :items="heartbeats",
        :loading="pending",
        :headers="headers",
        item-key="_id",
        select-all,
        expand
      )
        template(slot="progress")
          v-fade-transition
            v-progress-linear(height="2", indeterminate, color="primary")
        template(slot="items", slot-scope="props")
          tr(@click="props.expanded = !props.expanded")
            td(@click.stop="")
              v-checkbox-functional(v-model="props.selected", primary, hide-details)
            td {{ props.item._id }}
            td {{ props.item.expected_interval }}
            td
              v-layout
                v-flex
                  v-btn.error--text(
                    v-if="hasDeleteAnyHeartbeatAccess",
                    icon,
                    small,
                    @click.stop="showDeleteHeartbeatModal(props.item._id)"
                  )
                    v-icon(color="error") delete
        template(slot="expand", slot-scope="props")
          .secondary.lighten-1
            v-layout.py-3.secondary.lighten-2(row)
              v-textarea.my-2.mx-4.pa-0(
                :value="props.item.pattern | json",
                readonly,
                auto-grow,
                outline,
                hide-details,
                dark
              )
</template>

<script>
import rightsTechnicalExploitationHeartbeatMixin from '@/mixins/rights/technical/exploitation/heartbeat';

export default {
  mixins: [
    rightsTechnicalExploitationHeartbeatMixin,
  ],
  props: {
    heartbeats: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    showDeleteHeartbeatModal: {
      type: Function,
      default: () => () => {},
    },
    showDeleteSelectedHeartbeatsModal: {
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
        { text: this.$t('heartbeat.table.fields.id'), value: '_id' },
        { text: this.$t('heartbeat.table.fields.expectedInterval'), value: 'expected_interval' },
        { text: this.$t('common.actionsLabel'), value: 'actions' },
      ];
    },
  },
  watch: {
    heartbeats() {
      this.selected = [];
    },
  },
};
</script>

