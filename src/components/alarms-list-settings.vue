<template lang="pug">
  v-navigation-drawer(permanent, right, app)
    v-toolbar(color="blue darken-4")
      v-list
        v-list-tile
          v-list-tile-title(class="title white--text text-xs-center") Alarms list settings
    v-divider
    v-list(expand, class="pt-0")
      v-list-group
        v-list-tile(slot="activator", active-class="activeHeader") Title
        v-container
          v-text-field(placeholder="Widget title")
      v-list-group
        v-list-tile(slot="activator") Default Sort column
        v-container
          v-text-field(placeholder="Column name")
          v-select(:items="sortChoices", value="ASC")
      v-list-group
        v-list-tile(slot="activator") Periodic refresh
        v-container
          v-layout
            v-flex
              v-switch(
                v-model="isPeriodicRefreshEnable",
                color="green darken-3",
                hide-details,
              )
            v-flex
              v-text-field(
                class='pt-0',
                hide-details,
                type="number",
                :disabled="!isPeriodicRefreshEnable")
      v-list-group
        v-list-tile(slot="activator") Default number of elements/page
        v-container
          v-text-field(
            placeholder="Elements per page",
            type="number"
          )
      v-list-group
        v-list-tile(slot="activator") Filter on Open/Resolve
        v-container
          v-layout
            v-checkbox(
              label="Open",
              v-model="openCheckbox",
              hide-details
            )
            v-checkbox(
              label="Resolve",
              v-model="resolveCheckbox",
              hide-details
            )
      v-list-group
        v-list-tile(slot="activator") Filters
        v-container
          v-select(label="Select a filter")
      v-list-group
        v-list-tile(slot="activator") Info popup
    v-btn(
      color="green darken-4 white--text",
      depressed,
      fixed,
      right
    ) Save
</template>

<script>
export default {
  name: 'AlarmsListSettings',
  data() {
    return {
      sortChoices: ['ASC', 'DESC'],
      isPeriodicRefreshEnable: false,
      openCheckbox: true,
      resolveCheckbox: false,
    };
  },
};
</script>

<style scoped>
  .activeHeader {
    background-color: blue;
  }
</style>
