<template lang="pug">
  v-list-group
    v-list-tile(slot="activator")
      .validation-header(:class="{ 'error--text': hasAnyError }") {{ $t('settings.statsSelect.title') }}
        .font-italic.caption.ml-1 ({{ $t('settings.statsSelect.required') }})
    v-container
      v-alert(:value="errors.has('stats')", type="error") {{ $t('settings.statsSelect.required') }}
      v-btn(@click="showAddStatModal") {{ $t('modals.addStat.title.add') }}
      v-list.secondary(dark)
        v-list-group(v-for="(stat, key) in stats", :key="key")
          v-list-tile(slot="activator")
            v-list-tile-content
              v-list-tile-title {{ key }}
            v-list-tile-action
              v-layout
                v-btn.primary.mx-1(@click.stop="showEditStatModal(key, stat)", fab, small, depressed)
                  v-icon edit
                v-btn.error(@click.stop="showDeleteStatModal(key)", fab, small, depressed)
                  v-icon delete
          v-list-tile
            v-list-tile-title {{ $t('common.stat') }}: {{ stat.stat }}
          v-list-tile
            v-list-tile-title {{ $t('common.trend') }}: {{ stat.trend }}
          v-list-tile
            v-list-tile-title {{ $t('common.parameters') }}: {{ stat.parameters }}
</template>

<script>
import { omit } from 'lodash';

import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';
import formMixin from '@/mixins/form';
import formValidatorErrorMixin from '@/mixins/form/validator-error';

export default {
  inject: ['$validator'],
  mixins: [modalMixin, formMixin, formValidatorErrorMixin],
  model: {
    prop: 'stats',
    event: 'input',
  },
  props: {
    stats: {
      type: Object,
      default: () => ({}),
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  watch: {
    stats(value) {
      if (this.required) {
        this.$validator.validate('stats', value);
      }
    },
  },
  created() {
    if (this.required) {
      this.$validator.attach({
        name: 'stats',
        rules: 'required',
        getter: () => Object.values(this.stats),
        context: () => this,
        vm: this,
      });
    }
  },
  methods: {
    showAddStatModal() {
      this.showModal({
        name: MODALS.addStat,
        config: {
          title: this.$t('modals.addStat.title.add'),
          action: (stat) => {
            const newStat = {
              ...omit(stat, ['title', 'parameters']),

              parameters: stat.stat.options.reduce((acc, option) => {
                acc[option] = stat.parameters[option];

                return acc;
              }, {}),
            };

            this.updateField(stat.title, newStat);
          },
        },
      });
    },

    showEditStatModal(statTitle, stat) {
      this.showModal({
        name: MODALS.addStat,
        config: {
          title: this.$t('modals.addStat.title.edit'),
          stat,
          statTitle,
          action: newStat => this.updateAndMoveField(statTitle, newStat.title, omit(newStat, ['title'])),
        },
      });
    },

    showDeleteStatModal(statTitle) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.removeField(statTitle),
        },
      });
    },
  },
};
</script>
