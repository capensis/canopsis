<template lang="pug" functional>
  .v-input.v-input--selection-controls.v-input--checkbox(:class="{ 'input-group--active': props.inputValue }")
    .v-input__control
      .v-input__slot
        .v-input--selection-controls__input(@click="listeners.change(!props.inputValue)")
          .v-input--selection-controls__ripple.accent--text(v-ripple="{ center: true }")
          i.v-icon.material-icons(v-if="props.inputValue", :class="{ 'accent--text': props.inputValue }") check_box
          i.v-icon.material-icons(v-else, :class="{ 'accent--text': props.inputValue }") check_box_outline_blank
        label(v-show="props.label !== ''") {{ props.label }}
</template>
