<template lang="pug">
  modal-wrapper
    template(slot="title")
      span {{ $t('modals.eventFilterRule.editActions') }}
    template(slot="text")
      v-layout(justify-end)
        v-tooltip(top)
          v-btn(slot="activator", icon, @click="showCreateActionModal")
            v-icon.primary--text add
          span {{ $t('common.add') }}
      v-container
        h2 {{ $t('modals.eventFilterRule.actions') }}
        v-list(dark)
          draggable(v-model="actions")
            v-list-group(v-for="(action, index) in actions", :key="action.name")
              v-list-tile(slot="activator")
                v-list-tile-title {{index + 1}} - {{ action.type }} - {{ action.name || action.from }}
                v-btn(@click.stop="showEditActionModal(index)", icon)
                  v-icon(color="success") edit
                v-btn(@click.stop="deleteAction(index)", icon)
                  v-icon(color="error") delete
              v-list-tile
                v-layout(column)
                  div(v-if="action.name") {{ $t('common.name') }}: {{ action.name }}
                  div(v-if="action.value") {{ $t('common.value') }}: {{ action.value }}
                  div(v-if="action.description") {{ $t('common.description') }}: {{ action.description }}
                  div(v-if="action.from") {{ $t('common.from') }}: {{ action.from }}
                  div(v-if="action.to") {{ $t('common.to') }}: {{ action.to }}
    template(slot="actions")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import Draggable from 'vuedraggable';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.eventFilterRuleActions,
  $_veeValidate: {
    validator: 'new',
  },
  components: { Draggable, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    return {
      actions: [],
    };
  },
  mounted() {
    if (this.config.actions) {
      this.actions = [...this.config.actions];
    }
  },
  methods: {
    showCreateActionModal() {
      this.$modals.show({
        name: MODALS.eventFilterRuleCreateAction,
        config: {
          action: ruleAction => this.actions.push(ruleAction),
        },
      });
    },

    showEditActionModal(index) {
      this.$modals.show({
        name: MODALS.eventFilterRuleCreateAction,
        config: {
          ruleAction: this.actions[index],
          action: ruleAction => this.$set(this.actions, index, ruleAction),
        },
      });
    },

    deleteAction(index) {
      this.$delete(this.actions, index);
    },

    async submit() {
      await this.config.action(this.actions);
      this.$modals.hide();
    },
  },
};
</script>
