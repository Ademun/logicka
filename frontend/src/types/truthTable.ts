interface TruthTableVariable {
  Name: string;
  Value: boolean;
}

interface TruthTableEntry {
  Result: boolean;
  Variables: TruthTableVariable[];
}

export type {TruthTableEntry, TruthTableVariable}