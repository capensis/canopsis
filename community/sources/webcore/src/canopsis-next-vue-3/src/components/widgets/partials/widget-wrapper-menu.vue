<template>
  <v-menu>
    <template #activator="{ props }">
      <v-btn
        size="small"
        variant="text"
        icon
        v-bind="props"
      >
        <v-icon size="small">
          more_horiz
        </v-icon>
      </v-btn>
    </template>

    <v-list density="compact">
      <v-list-item
        v-if="hasSettingsSidebar"
        @click="showSettings"
      >
        <v-list-item-title>
          {{ $t('common.edit') }}
        </v-list-item-title>
      </v-list-item>

      <v-list-item @click="copyWidgetId">
        <v-list-item-title>
          {{ $t('view.copyWidgetId') }}
        </v-list-item-title>
      </v-list-item>

      <v-list-item @click="showDeleteWidgetModal">
        <v-list-item-title class="text-error">
          {{ $t('common.delete') }}
        </v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import { mapActions } from 'pinia';

import { MODALS, SIDE_BARS_BY_WIDGET_TYPES, WIDGET_TYPES } from '@/constants';

import { writeTextToClipboard } from '@/helpers/clipboard';

import { activeViewMixin } from '@/mixins/active-view';
import { useWidgetStore } from '@/stores/modules/widget';

export default {
  mixins: [
    activeViewMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    tab: {
      type: Object,
      required: true,
    },
  },
  computed: {
    hasSettingsSidebar() {
      return this.widget.type === WIDGET_TYPES.text;
    },
  },
  methods: {
    ...mapActions(useWidgetStore, {
      removeWidget: 'remove',
    }),

    showSettings() {
      this.$sidebar.show({
        name: SIDE_BARS_BY_WIDGET_TYPES[this.widget.type],
        config: {
          widget: {
            ...this.widget,
            tab: this.tab._id,
          },
        },
      });
    },

    /**
     * Show delete widget modal window
     */
    showDeleteWidgetModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeWidget({ id: this.widget._id });

            return this.fetchActiveView();
          },
        },
      });
    },

    /**
     * Copy widget id to clipboard.
     */
    async copyWidgetId() {
      try {
        await writeTextToClipboard(this.widget._id);

        this.$popups.success({ text: this.$t('success.widgetIdCopied') });
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: this.$t('errors.default') });
      }
    },
  },
};
</script>
