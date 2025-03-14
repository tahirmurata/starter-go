import type { SpawnOptions } from "bun";

const spawnOptions: SpawnOptions.OptionsObject = {
  stdin: "inherit",
  stdout: "inherit",
  stderr: "inherit",
};

const run = async () => {
  const templ = Bun.spawn(["bun", "run", "build:templ"], spawnOptions);
  const tailwindcss = Bun.spawn(
    ["bun", "run", "build:tailwindcss"],
    spawnOptions
  );
  const sqlc = Bun.spawn(["bun", "run", "build:sqlc"], spawnOptions);

  process.on("SIGINT", async () => {
    templ.kill;
    tailwindcss.kill;
    sqlc.kill;
    await Promise.all([templ.exited, tailwindcss.exited, sqlc.exited]);
  });

  await Promise.all([templ.exited, tailwindcss.exited, sqlc.exited]);

  const go = Bun.spawn(["bun", "run", "build:go"], spawnOptions);

  process.on("SIGINT", async () => {
    go.kill;
    await Promise.all([go.exited]);
  });

  await Promise.all([go.exited]);
};

run();
