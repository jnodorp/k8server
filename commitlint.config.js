module.exports = {
  extends: ["@commitlint/config-conventional"],
  ignores: [
    (message) => /((build\(go\))|(ci)): bump .+ from .+ to .+/.test(message),
  ],
};
