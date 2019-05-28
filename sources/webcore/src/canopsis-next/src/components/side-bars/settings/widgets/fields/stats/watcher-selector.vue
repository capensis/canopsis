<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsWatcherSelector.title') }}
    v-container(fluid)
      v-alert(:value="errors.has('entityId')", type="error") Watcher is required
      v-btn.primary(@click="showContextEntitySelectorModal") Select watcher
</template>

<script>
import { MODALS, CONTEXT_ENTITIES_TYPES } from '@/constants';

import modalMixin from '@/mixins/modal';
import formBaseMixin from '@/mixins/form/base';

import { getContextSearchByText } from '@/helpers/widget-search';

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
  created() {
    this.$validator.attach({
      name: 'entityId',
      rules: 'required:true',
      getter: () => this.value,
      context: () => this,
    });
  },
  methods: {
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
  },
};
</script>

