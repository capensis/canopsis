<template lang="pug">
  .ds-expand.ds-calendar-app
    v-navigation-drawer(
      v-model="drawer",
      fixed,
      app,
      :clipped="$vuetify.breakpoint.lgAndUp"
    )
      slot(name="drawerTop")
      slot(name="drawerPicker", :calendar="calendar", :picked="rebuild")
        .pa-3(v-if="calendar")
          ds-day-picker(:span="calendar.span", @picked="rebuild")
      slot(name="drawerBottom")

    v-toolbar.ds-app-calendar-toolbar(
      app,
      flat,
      fixed,
      color="white",
      :clipped-left="$vuetify.breakpoint.lgAndUp"
    )
      v-toolbar-title.ml-0(:style="toolbarStyle")
        v-toolbar-side-icon(@click.stop="drawer = !drawer")
        span.hidden-sm-and-down
        slot(name="title", :calendar="calendar")
      slot(
        name="today",
        :setToday="setToday",
        :todayDate="todayDate",
        :calendar="calendar"
      )
        v-tooltip(bottom)
          v-btn.ds-skinny-button(
            slot="activator",
            depressed,
            :icon="$vuetify.breakpoint.smAndDown",
            @click="setToday"
          )
            span(v-if="$vuetify.breakpoint.mdAndUp") {{ labels.today }}
            v-icon(v-else) {{ labels.todayIcon }}
          span {{ todayDate }}

      slot(
        name="prev",
        :prev="prev",
        :prevLabel="prevLabel",
        :calendar="calendar"
      )
        v-tooltip(bottom)
          v-btn.ds-light-forecolor.ds-skinny-button(
            slot="activator",
            icon,
            depressed,
            @click="prev"
          )
            v-icon keyboard_arrow_left
          span {{ prevLabel }}

      slot(
        name="next",
        :next="next",
        :nextLabel="nextLabel",
        :calendar="calendar"
      )
        v-tooltip(bottom)
          v-btn.ds-light-forecolor.ds-skinny-button(slot="activator", icon, depressed, @click="next")
            v-icon keyboard_arrow_right
          span {{ nextLabel }}

      slot(
        name="summary",
        :summary="summary",
        :calendar="calendar"
      )
        h1.title.ds-light-forecolor {{ summary }}
      v-spacer

      slot(
        name="view",
        :currentType="currentType",
        :types="types"
      )
        v-menu
          v-btn(flat, slot="activator") {{ currentType.label }}
            v-icon arrow_drop_down
          v-list
            v-list-tile(
              v-for="type in types",
              :key="type.id",
              @click="currentType = type"
            )
              v-list-tile-content
                v-list-tile-title {{ type.label }}
              v-list-tile-action {{ type.shortcut }}

      slot(name="menuRight")

    v-content.ds-expand
      v-container.ds-calendar-container(fluid, fill-height)
        ds-gestures(
          @swipeleft="next",
          @swiperight="prev"
        )
          div.ds-expand(v-if="currentType.schedule")
            slot(name="calendarAppAgenda", v-bind="{ $scopedSlots, $listeners, calendar, add, edit, viewDay }")
              ds-agenda(
                v-bind="{ $scopedSlots }",
                v-on="$listeners",
                :calendar="calendar",
                @add="add",
                @edit="edit",
                @view-day="viewDay"
              )
          div.ds-expand(v-else)
            slot(
              name="calendarAppCalendar",
              v-bind="{ $scopedSlots, $listeners, calendar, add, addAt, edit, viewDay, handleAdd, handleMove }"
            )
              ds-calendar(
                ref="calendar",
                v-bind="{ $scopedSlots }",
                v-on="$listeners",
                :calendar="calendar",
                @add="add",
                @add-at="addAt",
                @edit="edit",
                @view-day="viewDay",
                @added="handleAdd",
                @moved="handleMove",
                @resized="handleResize"
              )

        slot(name="calendarAppEventDialog", v-bind="{ $scopedSlots, $listeners, calendar, eventFinish }")
          ds-event-dialog(
            ref="eventDialog",
            v-bind="{ $scopedSlots }",
            v-on="$listeners",
            :calendar="calendar",
            @saved="eventFinish",
            @actioned="eventFinish"
          )

        slot(name="calendarAppOptions", v-bind="{ optionsVisible, optionsDialog, options, chooseOption }")
          v-dialog(
            ref="optionsDialog",
            v-model="optionsVisible",
            v-bind="optionsDialog",
            :fullscreen="$dayspan.fullscreenDialogs"
          )
            v-list
              template(v-for="option in options")
                v-list-tile(:key="option.text", @click="chooseOption(option)") {{ option.text }}

        slot(name="calendarAppPrompt", v-bind="{ promptVisible, promptDialog, promptQuestion, choosePrompt }")
          v-dialog(
            ref="promptDialog",
            v-model="promptVisible",
            v-bind="promptDialog"
          )
            v-card
              v-card-title {{ promptQuestion }}
              v-card-actions
                v-btn(color="primary", flat, @click="choosePrompt(true)") {{ labels.promptConfirm }}
                v-spacer
                v-btn(color="secondary", flat, @click="choosePrompt(false)") {{ labels.promptCancel }}

        slot(name="calendarAppAdd", v-bind="{ allowsAddToday, addToday }")
          v-fab-transition
            v-btn(
              class="ds-add-event-today",
              color="primary",
              fixed,
              bottom,
              right,
              fab,
              v-model="allowsAddToday",
              @click="addToday"
            )
              v-icon add

        slot(name="containerInside", v-bind="{ events, calendar }")
</template>

<script>
import { DsCalendarApp } from 'dayspan-vuetify/src/components';

export default {
  extends: DsCalendarApp,
  methods: {
    handleResize() {},
  },
};
</script>
