module.exports = {
    "extends": "airbnb",
    "env": {
        "browser": true,
    },
    "rules": {
        "import/no-extraneous-dependencies": [
            "error", {
                "devDependencies": true,
                "optionalDependencies": false,
                "peerDependencies": false
            }
        ]
    }
};
