import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

/**
 * @typedef {Object} TemplateTestingTest
 * @property {string} name
 * @property {string} description
 * @property {number} rule_type
 * @property {string} rule_name
 */

/**
 * Convert template testing test object to form
 *
 * @param {TemplateTestingTest} templateTestingTest
 * @returns {TemplateTestingTest}
 */
export const templateTestingTestToForm = (templateTestingTest = {}) => ({
  name: templateTestingTest.name ?? '',
  description: templateTestingTest.description ?? '',
  rule_type: templateTestingTest.rule_type ?? TEMPLATE_TESTING_TEST_TYPES.scenario,
  rule_name: templateTestingTest.rule_name ?? '',
});
