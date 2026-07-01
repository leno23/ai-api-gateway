"use client";

import { Line, LineChart, ResponsiveContainer } from "recharts";

type StatSparklineProps = {
  data: { t: string; v: number }[];
  color?: string;
};

export function StatSparkline({ data, color = "#007AFF" }: StatSparklineProps) {
  if (!data.length) return null;
  return (
    <div className="mt-2 h-10 w-full">
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={data}>
          <Line
            type="monotone"
            dataKey="v"
            stroke={color}
            strokeWidth={2}
            dot={false}
            isAnimationActive={false}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
