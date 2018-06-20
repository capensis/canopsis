<template lang="pug">
  v-navigation-drawer(
  :value="isPanelOpen",
  clipped,
  right,
  stateless,
  app,
  :temporary="$mq === 'mobile' || $mq === 'tablet'"
  )
    v-toolbar(color="blue darken-4")
      v-list
        v-list-tile
          v-list-tile-title.white--text.text-xs-center {{$t('settings.titles.alarmListSettings')}}
      v-icon.closeIcon(@click.stop="closePanel", color="white") close
    v-divider
    v-list.pt-0(expand)
      v-list-group
        v-list-tile(slot="activator", active-class="activeHeader") {{$t('common.title')}}
        v-container
          v-text-field(:placeholder="$t('settings.widgetTitle')")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{$t('settings.defaultSortColumn')}}
        v-container
          v-text-field(:placeholder="$t('settings.columnName')")
          v-select(:items="sortChoices", value="ASC")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{$t('settings.columnNames')}}
        v-container
          v-card
            v-layout.pt-2(justify-space-between)
              v-flex(xs3)
                v-layout.text-xs-center.pl-2(justify-space-between)
                  v-flex(xs1)
                    v-icon arrow_upward
                  v-flex(xs5)
                    v-icon arrow_downward
              v-flex.d-flex(xs3)
                div.text-xs-right.pr-2
                  v-icon(color="red") close
            v-layout(justify-center wrap)
              v-flex(xs11)
                v-text-field(:placeholder="$t('common.label')")
              v-flex(xs11)
                v-text-field(:placeholder="$t('common.value')")
          v-divider
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{$t('settings.periodicRefresh')}}
        v-container
          v-layout
            v-flex
              v-switch(
              v-model="isPeriodicRefreshEnable",
              color="green darken-3",
              hide-details,
              )
            v-flex
              v-text-field.pt-0(
                hide-details,
                type="number",
                :disabled="!isPeriodicRefreshEnable")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{$t('settings.defaultNumberOfElementsPerPage')}}
        v-container
          v-text-field(
          :placeholder="$t('settings.elementsPerPage')",
          type="number"
          )
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{$t('settings.filterOnOpenResolved')}}
        v-container
          v-layout
            v-checkbox(
            :label="$t('settings.open')",
            v-model="openCheckbox",
            hide-details
            )
            v-checkbox(
            :label="$t('settings.resolved')",
            v-model="resolveCheckbox",
            hide-details
            )
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{$t('settings.filters')}}
        v-container
          v-select(:label="$t('settings.selectAFilter')")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{$t('settings.infoPopup')}}
        info-popup-settings-item
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{$t('settings.moreInfosModal')}}
        v-container
          v-text-field(class="pa-0", textarea, hide-details)
    v-btn(
    color="green darken-4 white--text",
    depressed,
    fixed,
    right
    ) {{$t('common.save')}}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import InfoPopupSettingsItem from '@/components/other/info-popup/settings-item.vue';

const { mapGetters, mapActions } = createNamespacedHelpers('alarmsListSettings');

export default {
  components: {
    InfoPopupSettingsItem,
  },
  data() {
    return {
      sortChoices: ['ASC', 'DESC'],
      isPeriodicRefreshEnable: false,
      openCheckbox: true,
      resolveCheckbox: false,
    };
  },
  computed: {
    ...mapGetters(['isPanelOpen']),
  },
  methods: {
    ...mapActions(['closePanel']),
  },
};
</script>

<style scoped>
  .activeHeader {
    background-color: blue;
  }

  .closeIcon:hover {
    cursor: pointer;
  }
</style>
