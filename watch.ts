import type { SpawnOptions } from "bun";

const spawnOptions: SpawnOptions.OptionsObject = {
  stdin: "inherit",
  stdout: "inherit",
  stderr: "inherit",
};

const run = async () => {
  const go = Bun.spawn(["bun", "run", "watch:go"], spawnOptions);
  const templ = Bun.spawn(["bun", "run", "watch:templ"], spawnOptions);
  const tailwindcss = Bun.spawn(
    ["bun", "run", "watch:tailwindcss"],
    spawnOptions
  );
  const sqlc = Bun.spawn(["bun", "run", "watch:sqlc"], spawnOptions);

  process.on("SIGINT", async () => {
    go.kill;
    templ.kill;
    tailwindcss.kill;
    sqlc.kill;
    await Promise.all([
      go.exited,
      templ.exited,
      tailwindcss.exited,
      sqlc.exited,
    ]);
  });

  await Promise.all([go.exited, templ.exited, tailwindcss.exited, sqlc.exited]);
};

run();
