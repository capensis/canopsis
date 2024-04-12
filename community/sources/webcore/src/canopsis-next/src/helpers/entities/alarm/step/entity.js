import { ALARM_LIST_STEPS } from '@/constants';

/**
 * @typedef {
 *   'stateinc' |
 *   'statedec' |
 *   'changestate' |
 *   'statecounter' |
 *   'statusinc' |
 *   'statusdec' |
 *   'resolve' |
 *   'activate' |
 *   'ack' |
 *   'ackremove' |
 *   'pbhenter' |
 *   'pbhleave' |
 *   'assocticket' |
 *   'webhookstart' |
 *   'webhookprogress' |
 *   'webhookcomplete' |
 *   'webhookfail' |
 *   'declareticket' |
 *   'declareticketfail' |
 *   'declareticketruleprogress' |
 *   'declareticketrulecomplete' |
 *   'declareticketrulefailed' |
 *   'snooze' |
 *   'unsnooze' |
 *   'comment' |
 *   'metaalarmattach' |
 *   'metaalarmdetach' |
 *   'instructionstart' |
 *   'instructionpause' |
 *   'instructionresume' |
 *   'instructioncomplete' |
 *   'instructionabort' |
 *   'instructionfail' |
 *   'autoinstructionstart' |
 *   'autoinstructioncomplete' |
 *   'autoinstructionfail' |
 *   'junittestsuiteupdate' |
 *   'junittestcaseupdate'
 * } AlarmStepType
 */

/**
 * Checks if the provided type is a `statecounter` step type.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the type is `statecounter`, otherwise returns `false`.
 */
export const isStateCounterStepType = type => type === ALARM_LIST_STEPS.stateCounter;

/**
 * Checks if the provided type is a `resolve` step type.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the type is `resolve`, otherwise returns `false`.
 */
export const isResolveStepType = type => type === ALARM_LIST_STEPS.resolve;

/**
 * Checks if the provided type is a `activate` step type.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the type is `activate`, otherwise returns `false`.
 */
export const isActivateStepType = type => type === ALARM_LIST_STEPS.activate;

/**
 * Checks if the provided type is a `comment` step type.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the type is `comment`, otherwise returns `false`.
 */
export const isCommentStepType = type => type === ALARM_LIST_STEPS.comment;

/**
 * Checks if the provided type is a `assocticket` step type.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the type is `assocticket`, otherwise returns `false`.
 */
export const isAssocTicketStepType = type => type === ALARM_LIST_STEPS.assocTicket;

/**
 * Checks if the provided type is an acknowledgment step type.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the type is either `ack` or `ackremove`, otherwise returns `false`.
 */
export const isAckStepType = type => [
  ALARM_LIST_STEPS.ack,
  ALARM_LIST_STEPS.ackRemove,
].includes(type);

/**
 * Checks if the provided type is related to a snooze operation.
 *
 * This function determines whether the given `type` corresponds to either
 * a snooze or unsnooze operation within the alarm list steps.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the type is either `snooze` or `unsnooze`, otherwise returns `false`.
 */
export const isSnoozeStepType = type => [
  ALARM_LIST_STEPS.snooze,
  ALARM_LIST_STEPS.unsnooze,
].includes(type);

/**
 * Checks if the provided type is a status step type.
 *
 * This function determines whether the given type corresponds to either an increment or decrement
 * in status by checking if the type is included in the predefined status step types.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @return {boolean} Returns `true` if the type is either `statusinc` or `statusdec`, otherwise returns `false`.
 */
export const isChangeStatusStepType = type => [
  ALARM_LIST_STEPS.statusinc,
  ALARM_LIST_STEPS.statusdec,
].includes(type);

/**
 * Checks if the provided type is one of the change state step types.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the type is either `changestate`, `stateinc`, or `statedec`,
 * otherwise, returns `false`.
 */
export const isChangeStateStepType = type => [
  ALARM_LIST_STEPS.changeState,
  ALARM_LIST_STEPS.stateinc,
  ALARM_LIST_STEPS.statedec,
].includes(type);

/**
 * Checks if the provided type is related to JUnit step types.
 *
 * This function determines whether the given type is either a JUnit test suite update or a JUnit test case update.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the type is either `junittestsuiteupdate` or `junittestcaseupdate`,
 * otherwise returns `false`.
 */
export const isJunitStepType = type => [
  ALARM_LIST_STEPS.junitTestSuiteUpdate,
  ALARM_LIST_STEPS.junitTestCaseUpdate,
].includes(type);

/**
 * Checks if the provided type is a Pbehavior step type.
 *
 * This function determines if the given `type` corresponds to either entering or leaving a Pbehavior state,
 * based on predefined step types in `ALARM_LIST_STEPS`.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the `type` is either `pbhenter` or `pbhleave`, otherwise returns `false`.
 */
export const isPbehaviorStepType = type => [
  ALARM_LIST_STEPS.pbhenter,
  ALARM_LIST_STEPS.pbhleave,
].includes(type);

/**
 * Checks if the given type is one of the instruction step types.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the type is an instruction step type, otherwise `false`.
 */
export const isInstructionStepType = type => [
  ALARM_LIST_STEPS.instructionStart,
  ALARM_LIST_STEPS.instructionPause,
  ALARM_LIST_STEPS.instructionResume,
  ALARM_LIST_STEPS.instructionComplete,
  ALARM_LIST_STEPS.instructionAbort,
  ALARM_LIST_STEPS.instructionFail,
].includes(type);

/**
 * Checks if the provided type is one of the auto instruction step types.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the type is an auto instruction step type, otherwise `false`.
 */
export const isAutoInstructionStepType = type => [
  ALARM_LIST_STEPS.autoInstructionStart,
  ALARM_LIST_STEPS.autoInstructionComplete,
  ALARM_LIST_STEPS.autoInstructionFail,
].includes(type);

/**
 * Checks if the provided type is one of the declare ticket step types.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the type is one of the declare ticket step types, otherwise `false`.
 */
export const isDeclareTicketStepType = type => [
  ALARM_LIST_STEPS.declareTicket,
  ALARM_LIST_STEPS.declareTicketFail,
  ALARM_LIST_STEPS.declareTicketRuleInProgress,
  ALARM_LIST_STEPS.declareTicketRuleComplete,
  ALARM_LIST_STEPS.declareTicketRuleFail,
].includes(type);

/**
 * Checks if the provided type is a webhook step type.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the type is one of the webhook step types, otherwise `false`.
 */
export const isWebhookStepType = type => [
  ALARM_LIST_STEPS.webhookStart,
  ALARM_LIST_STEPS.webhookInProgress,
  ALARM_LIST_STEPS.webhookComplete,
  ALARM_LIST_STEPS.webhookFail,
].includes(type);

/**
 * Checks if the provided type is related to a meta alarm step.
 *
 * This function determines whether the given `type` corresponds to either
 * attaching or detaching a meta alarm by comparing it against predefined
 * step types in `ALARM_LIST_STEPS`.
 *
 * @param {AlarmStepType} type - The step type to check.
 * @returns {boolean} Returns `true` if the `type` is either `metaalarmattach` or `metaalarmdetach`,
 * otherwise returns `false`.
 */
export const isMetaAlarmStepType = type => [
  ALARM_LIST_STEPS.metaalarmattach,
  ALARM_LIST_STEPS.metaalarmdetach,
].includes(type);
