<template lang="pug">
  v-list-group(:disabled="pending")
    template(slot="appendIcon")
      v-progress-circular(v-show="pending", width="2", size="25", color="rgba(0,0,0,0.54)" indeterminate)
      v-icon(v-show="!pending") $vuetify.icons.expand
    template(slot="appendIcon")
    v-list-tile(slot="activator") {{ $t('settings.statsWatcherSelector.title') }}
    v-container(fluid)
      v-alert(:value="errors.has('entityId')", type="error") {{ $t('settings.statsWatcherSelector.required') }}
      div(v-if="watcher")
        v-list(two-line)
          v-list-tile(avatar, @click="")
            v-list-tile-content
              v-list-tile-title {{ watcher.name }}
              v-list-tile-sub-title {{ watcher._id }}
      v-btn.primary(@click="showContextEntitySelectorModal") {{ $t('common.select') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS, CONTEXT_ENTITIES_TYPES } from '@/constants';

import modalMixin from '@/mixins/modal';
import formBaseMixin from '@/mixins/form/base';

import { getContextSearchByText } from '@/helpers/widget-search';

const { mapActions } = createNamespacedHelpers('entity');

export default {
  inject: ['$validator'],
  mixins: [modalMixin, formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      pending: false,
      watcher: null,
    };
  },
  created() {
    this.$validator.attach({
      name: 'entityId',
      rules: 'required:true',
      getter: () => this.value,
      context: () => this,
    });
  },
  mounted() {
    this.fetchItem();
  },
  methods: {
    ...mapActions({
      fetchEntitiesListWithoutStore: 'fetchListWithoutStore',
    }),

    showContextEntitySelectorModal() {
      this.showModal({
        name: MODALS.contextEntitySelector,
        config: {
          entityId: this.value,
          filterPreparer: text => ({
            $and: [
              getContextSearchByText(text, ['_id']),

              { type: CONTEXT_ENTITIES_TYPES.watcher },
            ],
          }),
          action: entityId => this.updateModel(entityId),
        },
      });
    },

    async fetchItem() {
      if (this.value) {
        this.pending = true;

        const { entities } = await this.fetchEntitiesListWithoutStore({
          params: {
            _filter: { _id: this.value, type: CONTEXT_ENTITIES_TYPES.watcher },
          },
        });

        this.watcher = entities[0] || null;
        this.pending = false;
      }
    },
  },
};
</script>

