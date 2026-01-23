import { useEffect, useRef, useState } from "react";
import * as PIXI from "pixi.js";

type AgentUpdate = {
  id: number;
  x: number;
  z: number;
  rotation: number;
  type: string;
};
function getRotationFromVelocity(vx: number, vz: number) {
  return Math.atan2(vz, vx);
}
function App() {
  const [tick, setTick] = useState<number | null>(null);
  const stageRef = useRef<HTMLDivElement | null>(null);
  const appRef = useRef<PIXI.Application | null>(null);
  const spritesRef = useRef<Map<number, PIXI.Graphics>>(new Map());

  useEffect(() => {
    if (!stageRef.current) return;
    const app = new PIXI.Application({
      backgroundColor: 0x1e1e1e,
      resizeTo: stageRef.current,
      antialias: true,
      resolution: window.devicePixelRatio || 1,
    });
    stageRef.current.appendChild(app.view as HTMLCanvasElement);
    appRef.current = app;

    return () => {
      app.destroy(true, { children: true });
      appRef.current = null;
    };
  }, []);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws");
    ws.addEventListener("message", (ev) => {
      try {
        const data = JSON.parse(ev.data) as { tick: number; updated: AgentUpdate[] };
        setTick(data.tick);
        const app = appRef.current;
        if (!app) return;

        const W = app.renderer.width;
        const H = app.renderer.height;

        data.updated.forEach((u) => {
          if (u.type !== "agent") return;
          let g = spritesRef.current.get(u.id);
          const screenX = (u.x / 5) * W;
          const screenY = (u.z / 5) * H;

          if (!g) {
            g = new PIXI.Graphics();
            g.beginFill(0xffcc00);
            // g.drawRect(0, 0, 10, 10);
            const x1 = 10;
            const y1 = 10;
            const x2 = 5;
            const y2 = 20;
            const x3 = 15;
            const y3 = 20;

            g.moveTo(x1, y1);
            g.lineTo(x2, y2);
            g.lineTo(x3, y3);
            g.lineTo(x1, y1);
            g.endFill();
            // set pivot to the triangle center so rotation happens around the sprite center
            g.pivot.set(10, 15);
            g.x = screenX;
            g.y = screenY;
            g.rotation = u.rotation;
            g.zIndex = 1;
            app.stage.addChild(g);
            spritesRef.current.set(u.id, g);
          } else {
            // compute rotation from velocity; add +90deg because graphic's tip is upwards
            g.rotation = u.rotation;
            g.x += (screenX - g.x) * 0.6;
            g.y += (screenY - g.y) * 0.6;
          }
        });
      } catch (e) {
        // ignore
      }
    });
    return () => ws.close();
  }, []);

  return (
    <div style={{ display: "flex", gap: 12 }}>
      <div style={{ width: 500, height: 500 }} ref={stageRef} />
      <div style={{ width: 220 }}>
        <div style={{ marginBottom: 8 }}>
          <strong>Tick:</strong> {tick ?? "-"}
        </div>
        <div>
          <em>Agents rendered with Pixi.js</em>
        </div>
      </div>
    </div>
  );
}

export default App;
